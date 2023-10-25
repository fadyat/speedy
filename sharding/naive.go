package sharding

import (
	"slices"
)

type naive struct {

	// todo: rewrite to map
	//
	// shards is a slice of shards registered to the sharding algorithm.
	//
	// used slice for faster lookup and simplicity instead of map.
	shards []*Shard

	// hash is the hash function used to determine which shard a key belongs to.
	// Returned value will be modded by the number of shards.
	hash hashFn
}

func NewNaive(
	shards []*Shard,
	hashFn func(key string) uint32,
) Algorithm {
	return &naive{
		shards: shards,
		hash:   hashFn,
	}
}

func (s *naive) GetShard(key string) *Shard {
	idx := s.hash(key) % uint32(len(s.shards))
	return s.shards[idx]
}

func (s *naive) RegisterShard(shard *Shard) error {
	if s.exists(shard) != -1 {
		return ErrShardAlreadyRegistered
	}

	s.shards = append(s.shards, shard)
	return nil
}

func (s *naive) DeleteShard(shard *Shard) error {
	if idx := s.exists(shard); idx != -1 {
		s.shards = slices.Delete(s.shards, idx, idx+1)
		return nil
	}

	return ErrShardNotFound
}

func (s *naive) exists(shard *Shard) int {
	uk := shard.uniqueKey()

	return slices.IndexFunc(s.shards, func(s *Shard) bool {
		return s.uniqueKey() == uk
	})
}

func (s *naive) GetShards() []*Shard {
	return s.shards
}
