// Launches a single instance of a cache server.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fadyat/speedy/server"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func initLogger() {
	lg, _ := zap.NewProduction()
	zap.ReplaceGlobals(lg)
}

func main() {
	initLogger()

	// parse arguments
	port := flag.Int("port", 8080, "port number for gRPC server to listen on")
	configFile := flag.String("config", "", "filename of JSON config file with the info for initial nodes")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	// get new grpc id server
	grpcServer, cacheServer := server.NewGrpcServer(
		*configFile,
		server.DYNAMIC,
	)

	// run gRPC server
	zap.S().Infof("Running gRPC server on port %d...", *port)
	go grpcServer.Serve(listener)

	// register node with cluster
	cacheServer.RegisterNodeInternal()

	// run initial election
	cacheServer.RunElection()

	// start leader heartbeat monitor
	go cacheServer.StartLeaderHeartbeatMonitor(context.Background())

	// set up shutdown handler and block until sigint or sigterm received
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c

		zap.L().Info("Shutting down gRPC server...")
		grpcServer.Stop()

		os.Exit(0)
	}()

	// block indefinitely
	select {}
}
