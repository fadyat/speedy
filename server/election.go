package server

import (
	"context"
	"github.com/fadyat/speedy/api"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	ELECTION_RUNNING    = true
	NO_ELECTION_RUNNING = false
	NO_LEADER           = "NO LEADER"
)

// RunElection Run an election using the Bully Algorithm (https://en.wikipedia.org/wiki/Bully_algorithm)
func (s *CacheServer) RunElection() {
	// an individual node should run a single election process, not multiple concurrent ones
	s.electionLock.Lock()
	if s.electionStatus {
		zap.L().Info("Election already running, waiting for completion...")
		return
	}

	// update status to election running
	s.electionStatus = ELECTION_RUNNING
	s.electionLock.Unlock()

	// check status of every node
	localPID := int32(os.Getpid())
	zap.S().Infof("Running election. Local PID: %d", localPID)

	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}
		func() {
			// new identity service client
			ctx, cancel := context.WithTimeout(context.Background(), 3*Timeout)
			defer cancel()

			// make status request rpc
			c, err := s.NewCacheClient(node.Host, int(node.Port))
			if err != nil {
				zap.S().Infof("error creating grpc client to node node %s: %v", node.Id, err)
			}
			res, err := c.GetPid(ctx, &api.PidRequest{CallerPid: localPID})
			if err != nil {
				zap.S().Infof("PID request to node %s failed", node.Id)
				return
			}

			// if response has a higher PID (use node id as tie-breaker), we send it an election request and wait to receive
			// the election winner announcement.
			zap.S().Infof("Received PID %d from node %s (vs local PID %d on node %s)", res.Pid, node.Id, localPID, s.nodeID)
			if (localPID < res.Pid) || (res.Pid == localPID && s.nodeID < node.Id) {

				zap.S().Infof("Sending election request to node %s", node.Id)

				c, err := s.NewCacheClient(node.Host, int(node.Port))
				if err != nil {
					zap.S().Infof("error creating grpc client to node node %s: %v", node.Id, err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 3*Timeout)
				defer cancel()

				_, err = c.RequestElection(ctx, &api.ElectionRequest{CallerPid: localPID, CallerNodeId: s.nodeID})
				if err != nil {
					zap.S().Infof("Error requesting node %s run an election: %v", node.Id, err)
				}

				zap.L().Info("Waiting for decision...")
				// if after 5 seconds we receive no winner announcement, start the election process over
				select {
				case winner := <-s.decisionChan:
					if winner != "" {
						s.leaderID = winner
						zap.S().Infof("Received decision: Leader is node %s", s.leaderID)

						s.electionLock.Lock()
						s.electionStatus = NO_ELECTION_RUNNING
						s.electionLock.Unlock()
						return
					}
				case <-time.After(5 * Timeout):
					zap.L().Info("Timed out waiting for decision. Starting new election.")
					s.RunElection()

					s.electionLock.Lock()
					s.electionStatus = NO_ELECTION_RUNNING
					s.electionLock.Unlock()
					return
				}
			}
		}()
	}
	// if no other nodes have a higher PID, we are the winner
	s.leaderID = s.nodeID
	zap.S().Infof("set leader as self: %s", s.nodeID)

	// announce ourselves as winner to other nodes
	s.AnnounceNewLeader(s.leaderID)

	// reset election status
	s.electionLock.Lock()
	s.electionStatus = NO_ELECTION_RUNNING
	s.electionLock.Unlock()
}

// AnnounceNewLeader Announce new leader to all nodes
func (s *CacheServer) AnnounceNewLeader(winner string) {
	zap.S().Infof("Announcing node %s won election", winner)

	// if no response from any higher node IDs, declare self the winner and announce to all
	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*Timeout)

			// make status request rpc
			c, err := s.NewCacheClient(node.Host, int(node.Port))
			if err != nil {
				zap.S().Infof("error creating grpc client to node node %s: %v", node.Id, err)
			}

			_, err = c.UpdateLeader(ctx, &api.NewLeaderAnnouncement{LeaderId: winner})
			if err != nil {
				zap.S().Infof("Election winner announcement to node %s error: %v", node.Id, err)
			}
			defer cancel()
		}()
	}
}

// GetLeader Returns current leader
func (s *CacheServer) GetLeader(ctx context.Context) (*api.LeaderResponse, error) {
	// while there is no leader, run election
	ticker := time.NewTicker(3 * time.Second)

	for {
		if s.leaderID != NO_LEADER {
			return &api.LeaderResponse{Id: s.leaderID}, nil
		}
		s.RunElection()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			// if no leader was elected, wait 3 seconds then run another election
			if s.leaderID == NO_LEADER {
				zap.L().Info("No leader elected, waiting 3 seconds before trying again...")
			} else {
				return &api.LeaderResponse{Id: s.leaderID}, nil
			}
		}
	}
}

// StartLeaderHeartbeatMonitor Checks if leader is alive every 1 second. If no response for 3 seconds, new election is held.
func (s *CacheServer) StartLeaderHeartbeatMonitor(ctx context.Context) {
	// wait for decision to get leader
	zap.L().Info("Leader heartbeat monitor starting...")

	ticker := time.NewTicker(Timeout)
	for {
		// run heartbeat check every 1 second
		<-ticker.C

		// case 1: we are a follower
		if s.leaderID != s.nodeID {
			if !s.IsLeaderAlive() {
				zap.L().Info("Leader heartbeat failed, running new election")
				s.RunElection()
				zap.L().Info("Election done, leader heartbeat continuing")
			}
			select {
			case <-ctx.Done():
				return
			case <-s.shutdownChan:
				zap.L().Info("Received shutdown signal")
				break
			case <-time.After(Timeout):
				continue
			}
			// case 2: we are the leader, so check for any dead nodes and remove them from cluster
		} else {
			s.SyncNodesConfig()
		}
	}
}

func (s *CacheServer) SyncNodesConfig() {
	modified := false
	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}
		func() {

			// new identity service client
			c, err := s.NewCacheClient(node.Host, int(node.Port))
			if err != nil {
				zap.S().Infof("error creating grpc client to node node %s: %v", node.Id, err)
				delete(s.nodesConfig.Nodes, node.Id)
				modified = true
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*Timeout)
			defer cancel()

			zap.S().Infof("Checking health of node %s", node.Id)
			_, err = c.GetHeartbeat(ctx, &api.HeartbeatRequest{CallerNodeId: s.nodeID})
			if err != nil {
				zap.S().Infof("Node %s healthcheck returned error, removing from cluster", node.Id)
				delete(s.nodesConfig.Nodes, node.Id)
				modified = true
			}
		}()
	}

	// if cluster was modified, send out updated cluster config to other nodes
	if modified {
		zap.L().Info("Detected node config change, sending update to other nodes")
		s.updateClusterConfigInternal()
	}
}

// IsLeaderAlive Check if leader node is alive (3 second timeout)
func (s *CacheServer) IsLeaderAlive() bool {
	// make sure leader exists
	if s.leaderID == NO_LEADER {
		zap.S().Infof("IsLeaderAlive found leader doesn't exist")
		return false
	}
	// if this node is the leader, return true
	if s.nodeID == s.leaderID {
		return true
	}
	zap.S().Infof("leader is %s", s.leaderID)
	leader, ok := s.nodesConfig.Nodes[s.leaderID]
	if !ok {
		zap.S().Infof("leader %s does not exist", s.leaderID)
		return true
	}

	// new identity service client
	ctx, cancel := context.WithTimeout(context.Background(), 3*Timeout)
	defer cancel()

	// make status request rpc
	c, err := s.NewCacheClient(leader.Host, int(leader.Port))
	if err != nil {
		zap.S().Infof("error creating grpc client to node %s: %v", leader.Id, err)
		return false
	}

	_, err = c.GetHeartbeat(ctx, &api.HeartbeatRequest{CallerNodeId: s.nodeID})
	if err != nil {
		zap.S().Infof("Leader healthcheck returned error: %v", err)
		return false
	}
	return true
}

// UpdateLeader gRPC handler for updating the leader after
func (s *CacheServer) UpdateLeader(ctx context.Context, request *api.NewLeaderAnnouncement) (*api.GenericResponse, error) {
	zap.S().Infof("Received announcement leader is %s", request.LeaderId)
	s.leaderID = request.LeaderId
	s.decisionChan <- s.leaderID
	return &api.GenericResponse{Data: SUCCESS}, nil
}

// GetHeartbeat Return current status of this node (leader/follower)
func (s *CacheServer) GetHeartbeat(ctx context.Context, request *api.HeartbeatRequest) (*empty.Empty, error) {
	zap.S().Infof("Node %s returning heartbeat to node %s", s.nodeID, request.CallerNodeId)
	return &empty.Empty{}, nil
}

// GetPid gRPC handler that receives a request with the caller's PID and returns its own PID.
// If the PID is higher than the caller PID, we take over the election process.
func (s *CacheServer) GetPid(ctx context.Context, request *api.PidRequest) (*api.PidResponse, error) {
	local_pid := int32(os.Getpid())
	return &api.PidResponse{Pid: local_pid}, nil
}

// RequestElection gRPC handler which allows other nodes to ask this node to start a new election
func (s *CacheServer) RequestElection(ctx context.Context, request *api.ElectionRequest) (*api.GenericResponse, error) {
	// asynchronously run election and return successful response
	zap.S().Infof("received request for election from %s", request.CallerNodeId)
	go s.RunElection()
	return &api.GenericResponse{Data: SUCCESS}, nil
}
