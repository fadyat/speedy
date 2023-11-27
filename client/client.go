package client

import (
	"context"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/node"
	"github.com/fadyat/speedy/sharding"
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

	// fixme: current problem if we updating the nodes config, the sharding
	//  algorithm will not be updated, and we will have a difference.
	//  we need to update the sharding algorithm, when the nodes config
	//  is updated.

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
	n := c.nodesConfig.GetNode(c.algo.GetShard(key).ID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := n.Request().Get(ctx, &api.GetRequest{Key: key})
	if err != nil {
		return "", asClientError(err)
	}

	return resp.Value, nil
}

func (c *client) Put(key, value string) error {
	n := c.nodesConfig.GetNode(c.algo.GetShard(key).ID)

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
				if err := c.nodesConfig.Sync(); err != nil {
					errCh <- err
				}
			}
		}
	}()

	return errCh
}
