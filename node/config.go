package node

import (
	"fmt"
	"github.com/fadyat/speedy/sharding"
	"gopkg.in/yaml.v3"
	"os"
)

// NodesConfig is used as current state of the system, it is used to
// initialize the system, to update the nodes information, and to
// retrieve the current state of the system.
type NodesConfig struct {

	// Node.ID is used as a unique identifier for the node, and it is
	// used as a key in the Nodes map.
	Nodes map[string]*Node `yaml:"nodes"`
}

type NodesConfigOption func(*NodesConfig) error

func NewNodesConfig(
	opts ...NodesConfigOption,
) (*NodesConfig, error) {
	c := &NodesConfig{
		Nodes: make(map[string]*Node),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithInitialState(path string) NodesConfigOption {
	return func(c *NodesConfig) error {
		inode, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read initial state file: %w", err)
		}

		if err = yaml.Unmarshal(inode, c); err != nil {
			return fmt.Errorf("failed to unmarshal initial state file: %w", err)
		}

		for _, n := range c.Nodes {
			if e := n.RefreshClient(); e != nil {
				return fmt.Errorf("failed to refresh client: %w", e)
			}
		}

		return nil
	}
}

func (c *NodesConfig) GetShards() []*sharding.Shard {
	shards := make([]*sharding.Shard, 0, len(c.Nodes))
	for _, n := range c.Nodes {
		shards = append(shards, n.ToShard())
	}

	return shards
}
