package server

import (
	"context"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/eviction"
	"github.com/fadyat/speedy/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	KeyNotFoundMsg = "key not found"
)

type CacheServer struct {
	api.UnimplementedCacheServiceServer

	configPath string
	cache      eviction.Algorithm
}

func (s *CacheServer) Get(_ context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	if val, ok := s.cache.Get(req.Key); ok {
		return &api.GetResponse{Value: val}, nil
	}

	return nil, status.Error(codes.NotFound, KeyNotFoundMsg)
}

func (s *CacheServer) Put(_ context.Context, req *api.PutRequest) (*emptypb.Empty, error) {
	s.cache.Put(req.Key, req.Value)
	return &emptypb.Empty{}, nil
}

func (s *CacheServer) Len(_ context.Context, _ *emptypb.Empty) (*api.LengthResponse, error) {
	return &api.LengthResponse{Length: s.cache.Len()}, nil
}

func (s *CacheServer) GetClusterConfig(_ context.Context, _ *emptypb.Empty) (*api.ClusterConfig, error) {
	locallyStored, err := getLocallyStoredClusterConfig(s.configPath)
	if err != nil {
		return nil, err
	}

	return &api.ClusterConfig{Nodes: locallyStored}, nil
}

func getLocallyStoredClusterConfig(path string) ([]*api.Node, error) {
	// todo: can use cache here, to avoid reading from file system every time.
	//  and read only when the file is updated + update the cache.

	inode, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read initial state file: %w", err)
	}

	// todo: rewrite this peace of shit
	var nodes = struct {
		Nodes map[string]*node.Node `yaml:"nodes"`
	}{}

	if err = yaml.Unmarshal(inode, &nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial state file: %w", err)
	}

	var apiStyle = make([]*api.Node, 0, len(nodes.Nodes))
	for _, v := range nodes.Nodes {
		apiStyle = append(apiStyle, &api.Node{
			Id:   v.ID,
			Host: v.Host,
			Port: uint32(v.Port),
		})
	}

	return apiStyle, nil
}

func NewCacheServer(
	configPath string,
	algo eviction.Algorithm,
) *CacheServer {
	return &CacheServer{
		configPath: configPath,
		cache:      algo,
	}
}
