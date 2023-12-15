package main

import (
	"context"
	"errors"
	"fmt"
	. "github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/client"
	"github.com/fadyat/speedy/eviction"
	"github.com/fadyat/speedy/server"
	"github.com/fadyat/speedy/sharding"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type wg struct {
	sync.WaitGroup
}

func (w *wg) Go(f func()) {
	w.Add(1)

	go func() {
		defer w.Done()
		f()
	}()
}

const (
	clientConfigPath    = "client-init.yaml"
	serverConfigPath    = "server.yaml"
	serverNewConfigPath = "server-new.yaml"

	serverUsedConfigPath = "server-copy.yaml"

	defaultCacheCapacity = 1000
	defaultServerPort    = 8082
)

func startServer(ctx context.Context, w *wg, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		zap.S().Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	cacheServer := server.NewCacheServer(
		serverUsedConfigPath,
		eviction.NewLRU(defaultCacheCapacity),
	)
	RegisterCacheServiceServer(s, cacheServer)

	w.Go(func() {
		zap.S().Infof("starting server on port %d", port)
		switch e := s.Serve(lis); {
		case e == nil, errors.Is(e, grpc.ErrServerStopped):
			return
		default:
			zap.S().Fatalf("failed to start server: %v", e)
		}
	})

	w.Go(func() {
		<-ctx.Done()
		s.GracefulStop()
	})
}

func startClient(ctx context.Context, w *wg) {
	c, err := client.NewClient(
		clientConfigPath,
		sharding.ConsistentAlgorithm,
		client.WithSyncPeriod(2*time.Second),
	)
	if err != nil {
		zap.S().Fatalf("failed to create client: %v", err)
	}

	w.Go(func() {
		zap.S().Infof("starting client sync loop")
		for e := range c.SyncClusterConfig(ctx) {
			zap.S().Errorf("failed to sync cluster config: %v", e)
		}
	})
}

func replaceClusterConfig(ctx context.Context, periodicity time.Duration) {
	var (
		states  = []string{serverConfigPath, serverNewConfigPath}
		current = 0
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(periodicity):
			zap.S().Infof("replacing cluster config to %s", states[current])
			if err := copyFile(states[current], serverUsedConfigPath); err != nil {
				zap.S().Fatalf("failed to copy file: %v", err)
			}

			current = (current + 1) % len(states)
		}
	}
}

func initLogger() {
	lg, _ := zap.NewProduction()
	zap.ReplaceGlobals(lg)
}

func main() {
	initLogger()

	cleanup := withTemporaryFile(serverConfigPath, serverUsedConfigPath)
	zap.S().Infof("using config file: %s", serverUsedConfigPath)
	defer func() {
		cleanup()
		zap.S().Infof("removed config file: %s", serverUsedConfigPath)
	}()

	var (
		w            wg
		sigCtx, stop = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	)

	for i := 0; i < 3; i++ {
		port := defaultServerPort + i

		w.Go(func() {
			startServer(sigCtx, &w, port)
		})
	}

	w.Go(func() {
		startClient(sigCtx, &w)
	})

	w.Go(func() {
		replaceClusterConfig(sigCtx, 7*time.Second)
	})

	<-sigCtx.Done()
	zap.S().Info("received signal, stopping")
	stop()

	w.Wait()
}
