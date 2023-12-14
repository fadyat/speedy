package syncer

import (
	"context"
	"fmt"
)

type CacheConfig struct {
	Nodes      map[string]*NodeConfig `yaml:"nodes"`
	MasterInfo *NodeConfig            `yaml:"master"`
}

type Option func(*CacheConfig)

func WithMasterInfo(masterInfo *NodeConfig) Option {
	return func(c *CacheConfig) {
		c.MasterInfo = masterInfo
	}
}

func NewDefaultCacheConfig(opts ...Option) *CacheConfig {
	d := &CacheConfig{
		Nodes: make(map[string]*NodeConfig),
	}

	for _, o := range opts {
		o(d)
	}

	return d
}

type NodeConfig struct {
	ID   string `yaml:"id"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (n *NodeConfig) Same(other *NodeConfig) bool {
	if other == nil {
		return false
	}

	return n.ID == other.ID && n.Host == other.Host && n.Port == other.Port
}

func NewNodeConfig(id, host string, port int) *NodeConfig {
	return &NodeConfig{
		ID:   id,
		Host: host,
		Port: port,
	}
}

type stateFetcherFn func(context.Context) (*CacheConfig, error)
type syncStatesFn func(context.Context, *cacheConfigDiff) (bool, error)

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

	gotoState := diff(currentState, desiredState)
	return c.syncStates(ctx, gotoState)
}
