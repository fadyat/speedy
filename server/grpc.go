package server

import (
	"context"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/eviction"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	KeyNotFoundMsg = "key not found"
)

type CacheServer struct {
	api.UnimplementedCacheServiceServer

	cache eviction.Algorithm
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
	// current configuration need to be retrieved from the node file system.
	// updated via the leader node gRPC calls and updated in the node file system.
	//
	// right now, mocking the cluster config with the following:

	var locallyStored = []*api.Node{
		{
			Id:   "1",
			Host: "localhost",
			Port: 8081,
		},
		{
			Id:   "2",
			Host: "localhost",
			Port: 8082,
		},
		{
			Id:   "3",
			Host: "localhost",
			Port: 8083,
		},
	}

	return &api.ClusterConfig{Nodes: locallyStored}, nil
}

func NewCacheServer(
	algo eviction.Algorithm,
) *CacheServer {
	return &CacheServer{
		cache: algo,
	}
}
