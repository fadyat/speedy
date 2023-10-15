package server

import (
	"context"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/eviction"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CacheServer struct {
	api.UnimplementedCacheServiceServer

	cache eviction.Algorithm
}

func (s *CacheServer) Get(_ context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	if val, ok := s.cache.Get(req.Key); ok {
		return &api.GetResponse{Value: val}, nil
	}

	return &api.GetResponse{Value: eviction.KeyNotFoundMsg}, nil
}

func (s *CacheServer) Put(_ context.Context, req *api.PutRequest) (*emptypb.Empty, error) {
	s.cache.Put(req.Key, req.Value)
	return &emptypb.Empty{}, nil
}

func NewCacheServer(
	algo eviction.Algorithm,
) *CacheServer {
	return &CacheServer{
		cache: algo,
	}
}
