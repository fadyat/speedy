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
}

func NewClient(
	configPath string,
	algoType sharding.AlgorithmType,
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

	return &client{
		algo:        algo,
		nodesConfig: nodesConfig,
	}, nil
}

func (c *client) Get(key string) (string, error) {
	n := c.nodesConfig.Nodes[c.algo.GetShard(key).ID]

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := n.Request().Get(ctx, &api.GetRequest{Key: key})
	if err != nil {
		return "", asClientError(err)
	}

	return resp.Value, nil
}

func (c *client) Put(key, value string) error {
	n := c.nodesConfig.Nodes[c.algo.GetShard(key).ID]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := n.Request().Put(ctx, &api.PutRequest{Key: key, Value: value})
	if err != nil {
		return fmt.Errorf("failed to put value in cache: %w", err)
	}

	return nil
}
