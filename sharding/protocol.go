package sharding

import (
	"errors"
	"fmt"
)

type AlgorithmType string

const (
	NaiveAlgorithm      AlgorithmType = "naive"
	RendezvousAlgorithm AlgorithmType = "rendezvous"
	ConsistentAlgorithm AlgorithmType = "consistent"
)

func NewAlgo(
	algo AlgorithmType,
	shards []*Shard,
	hashFn func(key string) uint32,
) (Algorithm, error) {
	switch algo {
	case NaiveAlgorithm:
		return NewNaive(shards, hashFn), nil
	case RendezvousAlgorithm:
		return NewRendezvous(shards, hashFn), nil
	case ConsistentAlgorithm:
		return NewConsistent(shards, hashFn), nil
	default:
		return nil, fmt.Errorf("unknown sharding algorithm: %s", algo)
	}
}

var (
	ErrShardAlreadyRegistered = errors.New("shard already registered")
	ErrShardNotFound          = errors.New("shard not found")
)

type hashFn func(key string) uint32

type Shard struct {
	ID   string `yaml:"id"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (s *Shard) uniqueKey() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Algorithm interface {

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
