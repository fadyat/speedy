package sharding

import (
	"errors"
	"slices"
	"sort"
	"sync"
)

type consistent struct {
	mx          sync.RWMutex
	shards      map[uint32]*Shard
	orderedKeys []uint32
	hash        hashFn
}

func NewConsistent(
	shards []*Shard,
	hashFn func(key string) uint32,
) Algorithm {
	c := &consistent{
		hash:        hashFn,
		shards:      make(map[uint32]*Shard),
		orderedKeys: make([]uint32, 0, len(shards)),
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	for _, shard := range shards {
		if errors.Is(c.registerShardUnsafe(shard), ErrShardAlreadyRegistered) {
			continue
		}
	}

	slices.Sort(c.orderedKeys)
	return c
}

func (c *consistent) GetShard(key string) *Shard {
	var (
		hash    = c.hash(key)
		closest = sort.Search(len(c.orderedKeys), func(i int) bool {
			return c.orderedKeys[i] >= hash
		})
	)

	c.mx.RLock()
	defer c.mx.RUnlock()

	closest %= len(c.orderedKeys)
	return c.shards[c.orderedKeys[closest]]
}

func (c *consistent) RegisterShard(shard *Shard) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if err := c.registerShardUnsafe(shard); err != nil {
		return err
	}

	slices.Sort(c.orderedKeys)
	return nil
}

func (c *consistent) registerShardUnsafe(shard *Shard) error {
	var hash = c.hash("node" + shard.ID)
	if _, ok := c.shards[hash]; ok {
		return ErrShardAlreadyRegistered
	}

	c.shards[hash] = shard
	c.orderedKeys = append(c.orderedKeys, hash)
	return nil
}

func (c *consistent) DeleteShard(shard *Shard) error {
	var hash = c.hash("node" + shard.ID)

	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.shards[hash]; !ok {
		return ErrShardNotFound
	}

	delete(c.shards, hash)
	idx := slices.Index(c.orderedKeys, hash)
	c.orderedKeys = slices.Delete(c.orderedKeys, idx, idx+1)
	return nil
}

func (c *consistent) GetShards() []*Shard {
	var shards = make([]*Shard, 0, len(c.shards))

	c.mx.RLock()
	defer c.mx.RUnlock()

	for _, hash := range c.orderedKeys {
		shards = append(shards, c.shards[hash])
	}

	return shards
}
