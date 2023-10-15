package main

import (
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/eviction"
	"github.com/fadyat/speedy/server"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var (
	Version = "dev"
)

func main() {
	initLogger()

	c, err := NewConfig()
	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(
				grpcStyleLogger(zap.L()),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
		),
	)

	cacheServer := server.NewCacheServer(eviction.NewLRU(c.Cache.Capacity))
	api.RegisterCacheServiceServer(s, cacheServer)
	reflection.Register(s)

	listener, err := net.Listen("tcp", ":"+c.Server.GrpcPort)
	if err != nil {
		zap.L().Fatal("failed to create listener", zap.Error(err))
	}

	zap.L().Info("starting grpc server", zap.String("port", c.Server.GrpcPort), zap.String("version", Version))
	if e := s.Serve(listener); e != nil {
		zap.L().Fatal("failed to start grpc server", zap.Error(e))
	}
}
