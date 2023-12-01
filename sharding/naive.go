package sharding

import (
	"slices"
	"sync"
)

type naive struct {
	mx     sync.RWMutex
	shards map[string]*Shard
	keys   []string
	hash   hashFn
}

func NewNaive(
	shards []*Shard,
	hashFn func(key string) uint32,
) Algorithm {
	n := &naive{
		hash:   hashFn,
		shards: make(map[string]*Shard),
		keys:   make([]string, 0, len(shards)),
	}

	n.mx.Lock()
	defer n.mx.Unlock()

	for _, shard := range shards {
		logRegisterErr(n.registerShardUnsafe(shard))
	}

	return n
}

func (n *naive) GetShard(key string) *Shard {
	n.mx.RLock()
	defer n.mx.RUnlock()

	idx := n.hash(key) % uint32(len(n.keys))
	return n.shards[n.keys[idx]]
}

func (n *naive) RegisterShard(shard *Shard) error {
	n.mx.Lock()
	defer n.mx.Unlock()

	return n.registerShardUnsafe(shard)
}

func (n *naive) registerShardUnsafe(shard *Shard) error {
	if n.exists(shard) {
		return ErrShardAlreadyRegistered
	}

	n.shards[shard.ID] = shard
	n.keys = append(n.keys, shard.ID)
	return nil
}

func (n *naive) DeleteShard(shard *Shard) error {
	n.mx.Lock()
	defer n.mx.Unlock()

	if !n.exists(shard) {
		return ErrShardNotFound
	}

	delete(n.shards, shard.ID)
	idx := slices.Index(n.keys, shard.ID)
	n.keys = slices.Delete(n.keys, idx, idx+1)
	return nil
}

func (n *naive) exists(shard *Shard) bool {
	_, ok := n.shards[shard.ID]
	return ok
}

func (n *naive) GetShards() []*Shard {
	n.mx.RLock()
	defer n.mx.RUnlock()

	shards := make([]*Shard, 0, len(n.keys))
	for _, key := range n.keys {
		shards = append(shards, n.shards[key])
	}

	return shards
}
