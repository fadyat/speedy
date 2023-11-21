package client

import (
	"context"
	"errors"
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

func (c *client) Close() {
	for _, n := range c.nodesConfig.Nodes {
		if err := n.Close(); err != nil {
			zap.S().Errorf("failed to close node %s: %s", n.ID, err)
		}
	}
}

// todo: move this logic to node/config.go, find difference between two lists
//
//	mark the nodes with 3 states: add, remove, no change
//
// add status: need to create a new client, and register the shard
// remove status: need to delete the shard, and close the client
// no change status: do nothing
//
// by future design config need to be supported from the leader node,
// leader node will be responsible for updating the config, and the
// config will be propagated to all nodes.
//
// client periodically will fetch the config from the leader node,
// and update client-side configuration.
func (c *client) UpdateNodesConfig(config *node.NodesConfig) {
	if config == c.nodesConfig {
		return
	}

	for _, n := range c.nodesConfig.Nodes {
		c.nodesConfig.Nodes[n.ID].RemoveCandidate = true
	}

	for _, n := range config.Nodes {
		n.RemoveCandidate = false

		// if the node is already registered, then we don't need to
		// register it again, and we don't need to refresh the client.
		if c.nodesConfig.Nodes[n.ID] != nil {
			c.nodesConfig.Nodes[n.ID] = n
			continue
		}

		if err := n.RefreshClient(); err != nil {
			zap.S().Errorf("failed to refresh client: %s", err)
		}

		if err := c.algo.RegisterShard(&sharding.Shard{
			ID:   n.ID,
			Host: n.Host,
			Port: n.Port,
		}); err != nil && !errors.Is(err, sharding.ErrShardAlreadyRegistered) {
			zap.S().Errorf("failed to register shard: %s", err)
		}

		c.nodesConfig.Nodes[n.ID] = n
	}

	for _, n := range c.nodesConfig.Nodes {
		if n.RemoveCandidate {
			if err := c.algo.DeleteShard(&sharding.Shard{
				ID:   n.ID,
				Host: n.Host,
				Port: n.Port,
			}); err != nil && !errors.Is(err, sharding.ErrShardNotFound) {
				zap.S().Errorf("failed to delete shard: %s", err)
			}

			delete(c.nodesConfig.Nodes, n.ID)

			if err := n.Close(); err != nil {
				zap.S().Errorf("failed to closed connection to node %s: %s", n.ID, err)
			}
		}
	}

	c.nodesConfig = config
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
