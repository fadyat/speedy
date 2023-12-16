package server

import (
	"context"
	api "github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/node"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"time"
)

// RegisterNodeWithCluster gRPC handler for registering a new node with the cluster.
// New nodes call this RPC on the leader when they come online.
func (s *CacheServer) RegisterNodeWithCluster(ctx context.Context, nodeInfo *api.Node) (*api.GenericResponse, error) {
	// if we already have this node registered, return
	if _, ok := s.nodesConfig.Nodes[nodeInfo.Id]; ok {
		zap.S().Infof("Node %s already part of cluster", nodeInfo.Id)
		return &api.GenericResponse{Data: SUCCESS}, nil
	}

	// add node to hashmap config for easy lookup
	s.nodesConfig.Nodes[nodeInfo.Id] = node.NewNode(nodeInfo.Id, nodeInfo.Host, nodeInfo.Port)

	// send update to other nodes in cluster
	nodes := make([]*api.Node, 0)
	for _, node := range s.nodesConfig.Nodes {
		nodes = append(nodes, &api.Node{Id: node.Id, Host: node.Host, Port: node.Port})
	}
	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}
		func() {
			// create context
			reqCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			cfg := api.ClusterConfig{Nodes: nodes}
			c, err := s.NewCacheClient(node.Host, int(node.Port))
			if err != nil {
				zap.S().Error("unable to connect to node %s", node.Id)
				return
			}
			_, err = c.UpdateClusterConfig(reqCtx, &cfg)
			if err != nil {
				return
			}
		}()
	}
	return &api.GenericResponse{Data: SUCCESS}, nil
}

// UpdateClusterConfig gRPC handler for updating cluster config with incoming info
func (s *CacheServer) UpdateClusterConfig(ctx context.Context, req *api.ClusterConfig) (*empty.Empty, error) {
	zap.L().Info("Updating cluster config")
	s.nodesConfig.Nodes = make(map[string]*node.Node)
	for _, nodecfg := range req.Nodes {
		s.nodesConfig.Nodes[nodecfg.Id] = node.NewNode(nodecfg.Id, nodecfg.Host, nodecfg.Port)
	}
	return &empty.Empty{}, nil
}

// private function for server to send out updated cluster config to other nodes
func (s *CacheServer) updateClusterConfigInternal() {
	zap.L().Info("Sending out updated cluster config")

	// send update to other nodes in cluster
	nodes := make([]*api.Node, 0)
	for _, node := range s.nodesConfig.Nodes {
		nodes = append(nodes, &api.Node{Id: node.Id, Host: node.Host, Port: node.Port})
	}
	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}

		func() {
			// create context
			reqCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			cfg := api.ClusterConfig{Nodes: nodes}

			c, err := s.NewCacheClient(node.Host, int(node.Port))

			// skip node if error
			if err != nil {
				zap.S().Errorf("unable to connect to node %s", node.Id)
				return
			}

			_, err = c.UpdateClusterConfig(reqCtx, &cfg)
			if err != nil {
				zap.S().Infof("error sending cluster config to node %s: %v", node.Id, err)
			}
		}()
	}
}
