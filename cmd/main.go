package main

import (
	"github.com/fadyat/speedy/api"
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

	grpcServer, cacheServer := server.NewGrpcServer("", server.DYNAMIC)

	api.RegisterCacheServiceServer(s, cacheServer)
	reflection.Register(s)

	listener, err := net.Listen("tcp", ":"+c.Server.GrpcPort)
	if err != nil {
		zap.L().Fatal("failed to create listener", zap.Error(err))
	}

	zap.S().Infof("Running gRPC server on port %s...", c.Server.GrpcPort)
	go grpcServer.Serve(listener)

	// register node with cluster
	cacheServer.RegisterNodeInternal()

	// run initial election
	cacheServer.RunElection()

}
