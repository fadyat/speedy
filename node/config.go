package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/pkg"
	"github.com/fadyat/speedy/sharding"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"slices"
	"sync"
	"time"
)

const (

	// global timeout for gRPC calls.
	operationTimeout = 1 * time.Second
)

// NodesConfig is used as current state of the system, it is used to
// initialize the system, to update the nodes information, and to
// retrieve the current state of the system.
type NodesConfig struct {
	Nodes        Nodes `yaml:"nodes"`
	keys         []string
	mx           sync.RWMutex
	nodeSelector func(nc *NodesConfig) *Node
}

type NodesConfigOption func(*NodesConfig) error

func NewNodesConfig(
	opts ...NodesConfigOption,
) (*NodesConfig, error) {
	c := &NodesConfig{
		Nodes:        make(map[string]*Node),
		keys:         make([]string, 0),
		nodeSelector: oneAfterAnotherNodeSelector,
	}

	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *NodesConfig) GetNode(id string) *Node {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.Nodes[id]
}

func WithInitialState(path string) NodesConfigOption {
	return func(c *NodesConfig) error {
		cfg, err := pkg.FromYaml[NodesConfig](path)
		if err != nil {
			return fmt.Errorf("failed to setup nodes from config: %w", err)
		}

		c.Nodes = cfg.Nodes
		c.keys = c.Nodes.NodeIDs()
		for _, n := range c.Nodes {
			if e := n.RefreshClient(context.Background()); e != nil {
				return fmt.Errorf("failed to refresh client: %w", e)
			}
		}

		return nil
	}
}

func (c *NodesConfig) GetShards() []*sharding.Shard {
	c.mx.Lock()
	defer c.mx.Unlock()

	var shards = make([]*sharding.Shard, 0, len(c.Nodes))
	for _, k := range c.keys {
		shards = append(shards, c.Nodes[k].ToShard())
	}

	return shards
}

func (c *NodesConfig) Sync() (bool, error) {
	var sourceOfTruth = c.nodeSelector(c)
	if sourceOfTruth == nil {
		return false, errors.New("failed to select node")
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	desiredConfig, err := sourceOfTruth.Request().GetClusterConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return false, fmt.Errorf("failed to get cluster config: %w", err)
	}

	return c.syncStates(desiredConfig.Nodes)
}

func (c *NodesConfig) syncStates(desired []*api.Node) (bool, error) {
	var (
		wg    sync.WaitGroup
		errCh = make(chan error)
	)

	for _, d := range c.diff(desired) {
		wg.Add(1)

		go func(d *nodeDiff) {
			defer wg.Done()

			switch d.state {
			case nodeStateAdded:
				c.setupWithObservability(errCh, d)
			case nodeStateRemoved:
				c.teardownWithObservability(errCh, d)
			case nodeStateSynced:
				zap.S().Infof("node %s is synced", d.id)
			default:
				zap.S().Errorf("unknown node state: %d", d.state)
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	changed, err := collect(errCh)
	return changed, err
}

func (c *NodesConfig) diff(desired []*api.Node) map[string]*nodeDiff {
	var clientState = make(map[string]*nodeDiff)

	c.mx.RLock()
	for _, n := range c.Nodes {
		clientState[n.ID] = newNodeDiffFromNode(n, nodeStateRemoved)
	}
	c.mx.RUnlock()

	for _, n := range desired {
		if _, ok := clientState[n.Id]; !ok {
			clientState[n.Id] = newNodeDiffFromApiNode(n, nodeStateAdded)
			continue
		}

		clientState[n.Id].state = nodeStateSynced
	}

	return clientState
}

func (c *NodesConfig) setupNode(n *nodeDiff) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.Nodes[n.id]; ok {
		return fmt.Errorf("node %s already exists", n.id)
	}

	node := n.toNode()
	if e := node.RefreshClient(context.Background()); e != nil {
		return fmt.Errorf("failed to refresh client: %w", e)
	}

	c.Nodes[n.id] = node
	c.keys = append(c.keys, n.id)
	return nil
}

func (c *NodesConfig) setupWithObservability(
	errs chan<- error, d *nodeDiff,
) {
	zap.S().Infof("adding node %s", d.id)
	if err := c.setupNode(d); err != nil {
		errs <- fmt.Errorf("failed to setup node: %w", err)
		return
	}

	errs <- nil
}

func (c *NodesConfig) teardownNode(n *nodeDiff) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if node, ok := c.Nodes[n.id]; ok {
		delete(c.Nodes, n.id)
		idx := slices.Index(c.keys, n.id)
		c.keys = slices.Delete(c.keys, idx, idx+1)
		if e := node.Close(); e != nil {
			// ignoring the error, system state need to be updated any way
			zap.S().Errorf("failed to close node %s: %v", n.id, e)
		}
	}

	return nil
}

func (c *NodesConfig) teardownWithObservability(
	errs chan<- error, d *nodeDiff,
) {
	zap.S().Infof("removing node %s", d.id)
	if err := c.teardownNode(d); err != nil {
		errs <- fmt.Errorf("failed to teardown node: %w", err)
		return
	}

	errs <- nil
}

func collect(ch <-chan error) (bool, error) {
	var (
		errs    = make([]error, 0)
		changed = false
	)

	for err := range ch {
		if err != nil {
			errs = append(errs, err)
			continue
		}

		changed = true
	}

	if len(errs) > 0 {
		return changed, fmt.Errorf("failed to sync nodes: %w", errs[0])
	}

	return changed, nil
}
