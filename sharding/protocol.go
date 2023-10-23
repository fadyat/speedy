package sharding

import (
	"errors"
	"fmt"
)

var (
	ErrShardAlreadyRegistered = errors.New("shard already registered")
	ErrShardNotFound          = errors.New("shard not found")
)

type hashFn func(key string) uint32

type Shard struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (s *Shard) ConnStr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (s *Shard) uniqueKey() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type Sharding interface {

	// GetShard returns the shard that the key belongs to.
	GetShard(key string) *Shard

	// RegisterShard registers a shard to the sharding algorithm.
	RegisterShard(shard *Shard) error

	// DeleteShard deletes a shard from the sharding algorithm.
	DeleteShard(shard *Shard) error

	// GetShards returns all shards registered to the sharding algorithm.
	//
	// Currently, this method is only used for testing purposes.
	GetShards() []*Shard
}
