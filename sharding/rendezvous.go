package sharding

import "slices"

type rendezvous struct {
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
	if r.exists(shard) != -1 {
		return ErrShardAlreadyRegistered
	}

	r.shards = append(r.shards, shard)
	return nil
}

func (r *rendezvous) DeleteShard(shard *Shard) error {
	if idx := r.exists(shard); idx != -1 {
		r.shards = slices.Delete(r.shards, idx, idx+1)
		return nil
	}

	return ErrShardNotFound
}

func (r *rendezvous) exists(shard *Shard) int {
	uk := shard.uniqueKey()

	return slices.IndexFunc(r.shards, func(s *Shard) bool {
		return s.uniqueKey() == uk
	})
}

func (r *rendezvous) GetShards() []*Shard {
	return r.shards
}
