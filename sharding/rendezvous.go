package sharding

import (
	"slices"
	"sync"
)

type rendezvous struct {
	mx sync.RWMutex

	// todo: rewrite to map
	shards []*Shard
	hash   hashFn
}

func NewRendezvous(
	shards []*Shard,
	hashFn func(key string) uint32,
) Algorithm {
	return &rendezvous{
		shards: shards,
		hash:   hashFn,
	}
}

func (r *rendezvous) GetShard(key string) *Shard {
	r.mx.RLock()
	defer r.mx.RUnlock()

	idx := r.getMaxShardIdx(key)
	return r.shards[idx]
}

func (r *rendezvous) getMaxShardIdx(key string) uint32 {
	var maxIdx, maxHash uint32

	for i, shard := range r.shards {
		hash := r.hash(key + shard.ID)
		if hash > maxHash {
			maxHash, maxIdx = hash, uint32(i)
		}
	}

	return maxIdx
}

func (r *rendezvous) RegisterShard(shard *Shard) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if r.exists(shard) != -1 {
		return ErrShardAlreadyRegistered
	}

	r.shards = append(r.shards, shard)
	return nil
}

func (r *rendezvous) DeleteShard(shard *Shard) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if idx := r.exists(shard); idx != -1 {
		r.shards = slices.Delete(r.shards, idx, idx+1)
		return nil
	}

	return ErrShardNotFound
}

func (r *rendezvous) exists(shard *Shard) int {
	id := shard.ID

	return slices.IndexFunc(r.shards, func(s *Shard) bool {
		return s.ID == id
	})
}

func (r *rendezvous) GetShards() []*Shard {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return ascopy(r.shards)
}
