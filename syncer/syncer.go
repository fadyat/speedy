package syncer

import (
	"context"
	"fmt"
)

type CacheConfig struct {
	Nodes      map[string]*NodeConfig `yaml:"nodes"`
	MasterInfo *NodeConfig            `yaml:"master"`
}

func NewDefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Nodes: make(map[string]*NodeConfig),
	}
}

func (c *CacheConfig) NodesSlice() []*NodeConfig {
	var nodes = make([]*NodeConfig, 0, len(c.Nodes))
	for _, n := range c.Nodes {
		nodes = append(nodes, n)
	}

	return nodes
}

type NodeConfig struct {
	ID   string `yaml:"id"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewNodeConfig(id, host string, port int) *NodeConfig {
	return &NodeConfig{
		ID:   id,
		Host: host,
		Port: port,
	}
}

type stateFetcherFn func(context.Context) (*CacheConfig, error)
type syncStatesFn func(context.Context, map[string]*nodeDiff) (bool, error)

type Cluster struct {
	currentStateFetcher stateFetcherFn
	desiredStateFetcher stateFetcherFn
	syncStates          syncStatesFn
}

func NewCluster(
	currentStateFetcher, desiredStateFetcher stateFetcherFn,
	syncStates syncStatesFn,
) *Cluster {
	return &Cluster{
		currentStateFetcher: currentStateFetcher,
		desiredStateFetcher: desiredStateFetcher,
		syncStates:          syncStates,
	}
}

func (c *Cluster) Sync(ctx context.Context) (changed bool, err error) {
	currentState, err := c.currentStateFetcher(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get current state: %w", err)
	}

	desiredState, err := c.desiredStateFetcher(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get desired state: %w", err)
	}

	gotoState := diff(currentState.NodesSlice(), desiredState.NodesSlice())
	return c.syncStates(ctx, gotoState)
}
