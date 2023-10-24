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

func NewCacheServer(
	algo eviction.Algorithm,
) *CacheServer {
	return &CacheServer{
		cache: algo,
	}
}
