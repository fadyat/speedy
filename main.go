// Launches a single instance of a cache server.
package main

import (
	"flag"
	"fmt"
	"github.com/fadyat/speedy/server"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// parse arguments
	port := flag.Int("ports", 5005, "port number for gRPC server to listen on")
	capacity := flag.Int("capacity", 1000, "capacity of LRU cache")
	clientAuth := flag.Bool("enable_client_auth", true, "require client authentication (used for mTLS)")
	configFile := flag.String("config", "", "filename of JSON config file with the info for initial nodes")
	verbose := flag.Bool("verbose", false, "log events to terminal")

	flag.Parse()

	// set up listener TCP connectiion
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	// get new grpc id server
	grpcServer, cacheServer := server.NewCacheServer(
		*capacity,
		*configFile,
		*verbose,
		server.DYNAMIC,
		*clientAuth,
	)

	// run gRPC server
	cacheServer.LogInfoLevel(fmt.Sprintf("Running gRPC server on port %d...", *port))
	go grpcServer.Serve(listener)

	// register node with cluster
	cacheServer.RegisterNodeInternal()

	// run initial election
	cacheServer.RunElection()

	// start leader heartbeat monitor
	go cacheServer.StartLeaderHeartbeatMonitor()

	// set up shutdown handler and block until sigint or sigterm received
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c

		cacheServer.LogInfoLevel("Shutting down gRPC server...")
		grpcServer.Stop()

		//cacheServer.LogInfoLevel("Shutting down HTTP server...")
		//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		//defer cancel()
		//
		//if err := httpServer.Shutdown(ctx); err != nil {
		//	cacheServer.LogInfoLevel(fmt.Sprintf("Http server shutdown error: %s", err))
		//}
		os.Exit(0)
	}()

	// block indefinitely
	select {}
}
