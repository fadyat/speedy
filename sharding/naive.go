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
		if err := n.registerShardUnsafe(shard); err != nil {
			continue
		}
	}

	return n
}

func (s *naive) GetShard(key string) *Shard {
	s.mx.RLock()
	defer s.mx.RUnlock()

	idx := s.hash(key) % uint32(len(s.keys))
	return s.shards[s.keys[idx]]
}

func (s *naive) RegisterShard(shard *Shard) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.registerShardUnsafe(shard)
}

func (s *naive) registerShardUnsafe(shard *Shard) error {
	if s.exists(shard) {
		return ErrShardAlreadyRegistered
	}

	s.shards[shard.ID] = shard
	s.keys = append(s.keys, shard.ID)
	return nil
}

func (s *naive) DeleteShard(shard *Shard) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.exists(shard) {
		delete(s.shards, shard.ID)
		slices.DeleteFunc(s.keys, func(s string) bool {
			return s == shard.ID
		})
		return nil
	}

	return ErrShardNotFound
}

func (s *naive) exists(shard *Shard) bool {
	_, ok := s.shards[shard.ID]
	return ok
}

func (s *naive) GetShards() []*Shard {
	s.mx.RLock()
	defer s.mx.RUnlock()

	shards := make([]*Shard, 0, len(s.keys))
	for _, key := range s.keys {
		shards = append(shards, s.shards[key])
	}

	return shards
}
