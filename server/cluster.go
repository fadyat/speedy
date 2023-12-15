package server

import (
	"context"
	"fmt"
	. "github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/node"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// gRPC handler for registering a new node with the cluster.
// New nodes call this RPC on the leader when they come online.
func (s *CacheServer) RegisterNodeWithCluster(ctx context.Context, nodeInfo *Node) (*GenericResponse, error) {
	// if we already have this node registered, return
	if _, ok := s.nodesConfig.Nodes[nodeInfo.Id]; ok {
		s.logger.Infof("Node %s already part of cluster", nodeInfo.Id)
		return &GenericResponse{Data: SUCCESS}, nil
	}

	// add node to hashmap config for easy lookup
	s.nodesConfig.Nodes[nodeInfo.Id] = node.NewNode(nodeInfo.Id, nodeInfo.Host, nodeInfo.Port)

	// send update to other nodes in cluster
	var nodes []*Node
	for _, node := range s.nodesConfig.Nodes {
		nodes = append(nodes, &Node{Id: node.Id, Host: node.Host, Port: node.Port})
	}
	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}
		// create context
		reqCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		cfg := ClusterConfig{Nodes: nodes}
		c, err := s.NewCacheClient(node.Host, int(node.Port))
		if err != nil {
			s.logger.Errorf("unable to connect to node %s", node.Id)
			return nil, status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("Unable to connect to node being registered: %s", node.Id),
			)
		}
		c.UpdateClusterConfig(reqCtx, &cfg)
	}
	return &GenericResponse{Data: SUCCESS}, nil
}

// gRPC handler for updating cluster config with incoming info
func (s *CacheServer) UpdateClusterConfig(ctx context.Context, req *ClusterConfig) (*empty.Empty, error) {
	s.logger.Info("Updating cluster config")
	s.nodesConfig.Nodes = make(map[string]*node.Node)
	for _, nodecfg := range req.Nodes {
		s.nodesConfig.Nodes[nodecfg.Id] = node.NewNode(nodecfg.Id, nodecfg.Host, nodecfg.Port)
	}
	return &empty.Empty{}, nil
}

// private function for server to send out updated cluster config to other nodes
func (s *CacheServer) updateClusterConfigInternal() {
	s.logger.Info("Sending out updated cluster config")

	// send update to other nodes in cluster
	var nodes []*Node
	for _, node := range s.nodesConfig.Nodes {
		nodes = append(nodes, &Node{Id: node.Id, Host: node.Host, Port: node.Port})
	}
	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}
		// create context
		reqCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		cfg := ClusterConfig{Nodes: nodes}

		c, err := s.NewCacheClient(node.Host, int(node.Port))

		// skip node if error
		if err != nil {
			s.logger.Errorf("unable to connect to node %s", node.Id)
			continue
		}

		_, err = c.UpdateClusterConfig(reqCtx, &cfg)
		if err != nil {
			s.logger.Infof("error sending cluster config to node %s: %v", node.Id, err)
		}
	}
}
