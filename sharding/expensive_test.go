package sharding

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpensive_Flow(t *testing.T) {
	testcases := []struct {
		name   string
		shards []*Shard
		ops    func(s *expensive, fn hashFn)
	}{
		{
			name: "register shard",
			ops: func(s *expensive, _ hashFn) {
				require.NoError(t, s.RegisterShard(&Shard{Host: "localhost", Port: 8080}))
				require.Equal(t, 1, len(s.shards))
			},
		},
		{
			name: "register shard twice",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
			},
			ops: func(s *expensive, _ hashFn) {
				err := s.RegisterShard(&Shard{Host: "localhost", Port: 8080})
				require.Equal(t, ErrShardAlreadyRegistered, err)
			},
		},
		{
			name: "delete shard",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
			},
			ops: func(s *expensive, _ hashFn) {
				require.NoError(t, s.DeleteShard(&Shard{Host: "localhost", Port: 8080}))
				require.Equal(t, 0, len(s.shards))
			},
		},
		{
			name: "delete shard not found",
			ops: func(s *expensive, _ hashFn) {
				err := s.DeleteShard(&Shard{Host: "localhost", Port: 8080})
				require.Equal(t, ErrShardNotFound, err)
			},
		},
		{
			name: "get shard",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
				{Host: "localhost", Port: 8081},
			},
			ops: func(s *expensive, hash hashFn) {
				shard := s.GetShard("key")
				require.Equal(t, s.shards[hash("key")%uint32(len(s.shards))], shard)
			},
		},
		{
			name: "get shards",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
				{Host: "localhost:", Port: 8081},
			},
			ops: func(s *expensive, _ hashFn) {
				shards := s.GetShards()
				require.Equal(t, s.shards, shards)
			},
		},
		{
			name: "shard hash changes",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
				{Host: "localhost", Port: 8081},
			},
			ops: func(s *expensive, hash hashFn) {
				prev := s.GetShard("key")
				require.Equal(t, s.shards[hash("key")%uint32(len(s.shards))], prev)

				require.NoError(t, s.RegisterShard(&Shard{Host: "localhost:8082"}))
				curr := s.GetShard("key")
				require.NotEqual(t, prev, curr)
			},
		},
	}

	ln := func(s string) uint32 { return uint32(len(s)) }
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			sharding := NewExpensive(tc.shards, ln).(*expensive)
			tc.ops(sharding, ln)
		})
	}
}
