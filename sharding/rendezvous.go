package sharding

import (
	"slices"
	"sync"
)

type rendezvous struct {
	mx     sync.RWMutex
	shards map[string]*Shard
	keys   []string
	hash   hashFn
}

func NewRendezvous(
	shards []*Shard,
	hashFn func(key string) uint32,
) Algorithm {
	r := &rendezvous{
		hash:   hashFn,
		shards: make(map[string]*Shard),
		keys:   make([]string, 0, len(shards)),
	}

	r.mx.Lock()
	defer r.mx.Unlock()

	for _, shard := range shards {
		logRegisterErr(r.registerShardUnsafe(shard))
	}

	return r
}

func (r *rendezvous) GetShard(key string) *Shard {
	r.mx.RLock()
	defer r.mx.RUnlock()

	idx := r.getMaxShardIdx(key)
	return r.shards[r.keys[idx]]
}

func (r *rendezvous) getMaxShardIdx(key string) uint32 {
	var maxIdx, maxHash uint32

	for i, shardID := range r.keys {
		hash := r.hash(key + shardID)
		if hash > maxHash {
			maxHash, maxIdx = hash, uint32(i)
		}
	}

	return maxIdx
}

func (r *rendezvous) RegisterShard(shard *Shard) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	return r.registerShardUnsafe(shard)
}

func (r *rendezvous) registerShardUnsafe(shard *Shard) error {
	if r.exists(shard) {
		return ErrShardAlreadyRegistered
	}

	r.shards[shard.ID] = shard
	r.keys = append(r.keys, shard.ID)
	return nil
}

func (r *rendezvous) DeleteShard(shard *Shard) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if !r.exists(shard) {
		return ErrShardNotFound
	}

	delete(r.shards, shard.ID)
	idx := slices.Index(r.keys, shard.ID)
	r.keys = slices.Delete(r.keys, idx, idx+1)
	return nil
}

func (r *rendezvous) exists(shard *Shard) bool {
	id := shard.ID
	_, ok := r.shards[id]
	return ok
}

func (r *rendezvous) GetShards() []*Shard {
	r.mx.RLock()
	defer r.mx.RUnlock()

	shards := make([]*Shard, 0, len(r.keys))
	for _, key := range r.keys {
		shards = append(shards, r.shards[key])
	}

	return shards
}
