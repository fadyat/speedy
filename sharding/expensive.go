package sharding

import "slices"

type expensive struct {

	// shards is a slice of shards registered to the sharding algorithm.
	//
	// used slice for faster lookup and simplicity instead of map.
	shards []*Shard

	// hash is the hash function used to determine which shard a key belongs to.
	// Returned value will be modded by the number of shards.
	hash hashFn
}

func NewExpensive(
	shards []*Shard,
	hashFn func(key string) uint32,
) Sharding {
	return &expensive{
		shards: shards,
		hash:   hashFn,
	}
}

func (s *expensive) GetShard(key string) *Shard {
	idx := s.hash(key) % uint32(len(s.shards))
	return s.shards[idx]
}

func (s *expensive) RegisterShard(shard *Shard) error {
	if s.exists(shard) != -1 {
		return ErrShardAlreadyRegistered
	}

	s.shards = append(s.shards, shard)
	return nil
}

func (s *expensive) DeleteShard(shard *Shard) error {
	idx := s.exists(shard)
	if idx == -1 {
		return ErrShardNotFound
	}

	s.shards = slices.Delete(s.shards, idx, idx+1)
	return nil
}

func (s *expensive) exists(shard *Shard) int {
	uk := shard.uniqueKey()

	for i, sh := range s.shards {
		if sh.uniqueKey() == uk {
			return i
		}
	}

	return -1
}

func (s *expensive) GetShards() []*Shard {
	return s.shards
}
