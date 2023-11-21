package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/eviction"
	"github.com/fadyat/speedy/server"
	"github.com/fadyat/speedy/sharding"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	defaultCacheCapacity = 2000
	defaultServerPort    = 50051

	singleNodeConfig = `
nodes:
  1:
    id: 1
    host: localhost
    port: 50051
`

	multipleNodesConfig = `
nodes:
  1:
    id: 1
    host: localhost
    port: 50051
  2:
    id: 2
    host: localhost
    port: 50052
  3:
    id: 3
    host: localhost
    port: 50053`
)

func withTemporaryFile(t require.TestingT, content string) (path string, cleanup func()) {
	f, err := os.CreateTemp("", "speedy-client-test.yaml")
	require.NoError(t, err)

	_, err = f.WriteString(content)
	if err != nil {
		require.NoError(t, f.Close())
		require.NoError(t, os.Remove(f.Name()))
	}

	require.NoError(t, err)

	return f.Name(), func() {
		require.NoError(t, os.Remove(f.Name()))
	}
}

func upServer(ctx context.Context, wg *sync.WaitGroup, t require.TestingT, port int) error {
	s := grpc.NewServer()
	cacheServer := server.NewCacheServer(eviction.NewLRU(defaultCacheCapacity))
	api.RegisterCacheServiceServer(s, cacheServer)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	go func() {
		defer wg.Done()

		log.Printf("starting grpc server on port %d", port)
		switch err = s.Serve(l); {
		case err == nil || errors.Is(grpc.ErrServerStopped, err):
		default:
			require.NoError(t, err, "failed to start grpc server")
		}
	}()

	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()

	return nil
}

func TestClient_Flow(t *testing.T) {
	testcases := []struct {
		name   string
		pre    func(c Client)
		verify func(c Client)
	}{
		{
			name: "cache miss",
			pre:  func(c Client) {},
			verify: func(c Client) {
				_, e := c.Get("key")
				require.Equal(t, ErrCacheMiss, e)
			},
		},
		{
			name: "cache hit",
			pre: func(c Client) {
				require.NoError(t, c.Put("key", "value"))
			},
			verify: func(c Client) {
				v, e := c.Get("key")
				require.Equal(t, "value", v)
				require.NoError(t, e)
			},
		},
		{
			name: "cache 1000 keys",
			pre: func(c Client) {
				for i := 0; i < defaultCacheCapacity/2; i++ {
					require.NoError(t, c.Put(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)))
				}
			},
			verify: func(c Client) {
				for i := 0; i < defaultCacheCapacity/2; i++ {
					v, e := c.Get(fmt.Sprintf("key%d", i))
					require.Equal(t, fmt.Sprintf("value%d", i), v)
					require.NoError(t, e)
				}
			},
		},
		{
			name: "out of capacity",
			pre: func(c Client) {
				for i := 0; i < defaultCacheCapacity*2; i++ {
					require.NoError(t, c.Put(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)))
				}
			},
			verify: func(c Client) {
				var notFound int
				for i := 0; i < defaultCacheCapacity*2; i++ {
					_, e := c.Get(fmt.Sprintf("key%d", i))
					if errors.Is(e, ErrCacheMiss) {
						notFound++
					}
				}

				require.Equal(t, defaultCacheCapacity, notFound)
			},
		},
	}

	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)

	wg.Add(1)
	require.NoError(t, upServer(ctx, &wg, t, defaultServerPort))

	path, cleanup := withTemporaryFile(t, singleNodeConfig)
	defer cleanup()

	c, err := NewClient(path, sharding.RendezvousAlgorithm)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.pre(c)
			tc.verify(c)
		})
	}

	c.Close()
	cancel()
	wg.Wait()
}

func TestClient_MultipleNodes(t *testing.T) {
	var (
		nodes = 3
	)

	testcases := []struct {
		name   string
		pre    func(c Client)
		verify func(c Client)
	}{
		{
			name: "as in single node out of capacity",
			pre: func(c Client) {
				for i := 0; i < defaultCacheCapacity*2; i++ {
					require.NoError(t, c.Put(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)))
				}
			},
			verify: func(c Client) {
				for i := 0; i < defaultCacheCapacity*2; i++ {
					_, e := c.Get(fmt.Sprintf("key%d", i))
					require.NoError(t, e)
				}
			},
		},
	}

	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)

	for i := 0; i < nodes; i++ {
		wg.Add(1)
		require.NoError(t, upServer(ctx, &wg, t, defaultServerPort+i))
	}

	path, cleanup := withTemporaryFile(t, multipleNodesConfig)
	defer cleanup()

	c, err := NewClient(path, sharding.RendezvousAlgorithm)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.pre(c)
			tc.verify(c)
		})
	}

	c.Close()
	cancel()
	wg.Wait()
}

func BenchmarkClient_Get(b *testing.B) {
	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)

	wg.Add(1)
	require.NoError(b, upServer(ctx, &wg, b, defaultServerPort))

	path, cleanup := withTemporaryFile(b, multipleNodesConfig)
	defer cleanup()

	c, err := NewClient(path, sharding.RendezvousAlgorithm)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, e := c.Get(fmt.Sprintf("key%d", i))
		require.NoError(b, e)
	}

	cancel()
	wg.Wait()
}

// current benchmark is not representative, will be changed in the future
//
// benchmarking of sharding algorithms it's really hard to do
//
// need to make a design decision on how to benchmark the algorithms
// 1. benchmark the algorithms with a fixed number of keys and nodes
// 2. benchmark the algorithms with a fixed number of keys and a variable number of nodes
// 3. benchmark the algorithms with a variable number of keys and a fixed number of nodes
// 4. benchmark the algorithms with a variable number of keys and a variable number of nodes
//
// each benchmark need to represent cache misses/hits, how nodes are filled up, how keys are distributed
//
// the best one is to benchmark cases, when nodes are added/removed, how keys are redistributed
//
// the hardest part is to make the benchmark representative of the real world
func BenchmarkClient_70PercentsRead(b *testing.B) {
	testcases := []struct {
		name     string
		nodes    int
		algoType sharding.AlgorithmType
	}{
		{
			name:     "naive",
			algoType: sharding.NaiveAlgorithm,
		},
		{
			name:     "rendezvous",
			algoType: sharding.RendezvousAlgorithm,
		},
		{
			name:     "consistent",
			algoType: sharding.ConsistentAlgorithm,
		},
	}

	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
		nodes       = 3
		maxKeys     = defaultCacheCapacity / nodes
		readRate    = 0.7
	)

	for i := 0; i < nodes; i++ {
		wg.Add(1)
		require.NoError(b, upServer(ctx, &wg, b, defaultServerPort+i))
	}

	for _, tc := range testcases {
		b.Run(tc.name, func(b *testing.B) {
			path, cleanup := withTemporaryFile(b, multipleNodesConfig)
			defer cleanup()

			c, err := NewClient(path, tc.algoType)
			require.NoError(b, err)

			var (
				cacheMisses atomic.Int64
				cacheHits   atomic.Int64
			)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				key := fmt.Sprintf("key%d", i%maxKeys)
				if i%int(1/readRate) == 0 {
					_, e := c.Get(key)
					switch {
					case e == nil:
						cacheHits.Add(1)
					default:
						cacheMisses.Add(1)
					}
				}

				require.NoError(b, c.Put(key, key))
			}

			// waiting some time to free up the servers
			b.StopTimer()
			c.Close()
			<-time.After(3 * time.Second)
		})
	}

	cancel()
	wg.Wait()
}
