package client

import (
	"context"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/node"
	"github.com/fadyat/speedy/sharding"
	"go.uber.org/zap"
	"hash/crc32"
	"time"
)

type client struct {
	nodesConfig *node.NodesConfig
	algo        sharding.Algorithm

	syncPeriod time.Duration
	errChSize  int
}

type Option func(*client)

func WithSyncPeriod(period time.Duration) Option {
	return func(c *client) {
		c.syncPeriod = period
	}
}

func NewClient(
	configPath string,
	algoType sharding.AlgorithmType,
	opts ...Option,
) (Client, error) {
	nodesConfig, err := node.NewNodesConfig(node.WithInitialState(configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize nodes config: %w", err)
	}

	algo, err := sharding.NewAlgo(
		algoType,
		nodesConfig.GetShards(),
		func(k string) uint32 { return crc32.ChecksumIEEE([]byte(k)) },
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sharding algorithm: %w", err)
	}

	c := &client{
		nodesConfig: nodesConfig,
		algo:        algo,
		syncPeriod:  2 * time.Second,
		errChSize:   10,
	}

	for _, o := range opts {
		o(c)
	}

	return c, nil
}

func (c *client) Get(key string) (string, error) {
	shard := c.algo.GetShard(key)
	if shard == nil {
		return "", ErrCacheMiss
	}

	n := c.nodesConfig.GetNode(shard.ID)
	if n == nil {
		zap.L().Info("sharding and nodes config are not synced, got outdated shard")
		return "", ErrCacheMiss
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := n.Request().Get(ctx, &api.GetRequest{Key: key})
	if err != nil {
		return "", asClientError(err)
	}

	return resp.Value, nil
}

func (c *client) Put(key, value string) error {
	shard := c.algo.GetShard(key)
	if shard == nil {
		return ErrCacheMiss
	}

	n := c.nodesConfig.GetNode(shard.ID)
	if n == nil {
		zap.L().Info("sharding and nodes config are not synced, got outdated shard")
		return ErrCacheMiss
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := n.Request().Put(ctx, &api.PutRequest{Key: key, Value: value})
	if err != nil {
		return fmt.Errorf("failed to put value in cache: %w", err)
	}

	return nil
}

func (c *client) SyncClusterConfig(ctx context.Context) <-chan error {
	errCh := make(chan error, c.errChSize)

	go func() {
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(c.syncPeriod):
				// INFO: we don't guarantee consistency in sync process.
				//
				// in current implementation syncing between nodes config and shards
				//  are done in non-atomic way.
				//
				// we're increasing cache miss rate in this case, but our cache don't
				//  have high consistency guarantees, so it's ok.
				//
				// Example of inconsistency:
				// - we have 2 nodes, and 2 shards
				// - nodes config is updated, and we have 2 node, where on of them is
				//   different from the previous one (e.g. removed and new one added)
				//
				//  [1, 2] -> [1, 3], in that case sharding algorithm will return shard
				//   with id 2, but node with id 2 is not exists anymore, that means
				//   that she is down, and we will in any cases get ErrCacheMiss.
				//
				// - shards config is updated, nodes and shards config are synced now

				changed, err := c.nodesConfig.Sync()
				if err != nil {
					errCh <- fmt.Errorf("failed to sync nodes config: %w", err)
					continue
				}

				if !changed {
					zap.L().Debug("nodes config is not changed")
					continue
				}

				zap.L().Debug("nodes config is changed, syncing shards")
				sharding.SyncShards(c.algo, c.nodesConfig.GetShards())
			}
		}
	}()

	return errCh
}
