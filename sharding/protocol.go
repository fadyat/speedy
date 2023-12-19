package sharding

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
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
	Port uint32 `yaml:"port"`
}

type Algorithm interface {
	GetShard(key string) *Shard
	RegisterShard(shard *Shard) error
	DeleteShard(shard *Shard) error
	GetShards() []*Shard
}

func logRegisterErr(err error) {
	switch {
	case errors.Is(err, ErrShardAlreadyRegistered), err == nil:
	default:
		zap.L().Error("failed to register shard", zap.Error(err))
	}
}

func logDeleteErr(err error) {
	switch {
	case errors.Is(err, ErrShardNotFound), err == nil:
	default:
		zap.L().Error("failed to delete shard", zap.Error(err))
	}
}

func SyncShards(algo Algorithm, desired []*Shard) {
	var (
		desiredMap = make(map[string]*Shard, len(desired))
		existing   = algo.GetShards()
	)

	for _, shard := range desired {
		desiredMap[shard.ID] = shard
	}

	for _, shard := range existing {
		if _, ok := desiredMap[shard.ID]; !ok {
			logDeleteErr(algo.DeleteShard(shard))
		}
	}

	for _, shard := range desired {
		if err := algo.RegisterShard(shard); err != nil {
			logRegisterErr(err)
		}
	}
}
