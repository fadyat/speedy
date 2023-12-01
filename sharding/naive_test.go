package sharding

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNaive_Flow(t *testing.T) {
	testcases := []struct {
		name   string
		shards []*Shard
		ops    func(s *naive, fn hashFn)
	}{
		{
			name: "register shard",
			ops: func(s *naive, _ hashFn) {
				require.NoError(t, s.RegisterShard(&Shard{Host: "localhost", Port: 8080}))
				require.Equal(t, 1, len(s.shards))
			},
		},
		{
			name: "register shard twice",
			shards: []*Shard{
				{ID: "1", Host: "localhost", Port: 8080},
			},
			ops: func(s *naive, _ hashFn) {
				err := s.RegisterShard(&Shard{Host: "localhost", Port: 8080, ID: "1"})
				require.Equal(t, ErrShardAlreadyRegistered, err)
			},
		},
		{
			name: "delete shard",
			shards: []*Shard{
				{ID: "1", Host: "localhost", Port: 8080},
			},
			ops: func(s *naive, _ hashFn) {
				require.NoError(t, s.DeleteShard(&Shard{Host: "localhost", Port: 8080, ID: "1"}))
				require.Equal(t, 0, len(s.shards))
			},
		},
		{
			name: "delete shard not found",
			ops: func(s *naive, _ hashFn) {
				err := s.DeleteShard(&Shard{Host: "localhost", Port: 8080})
				require.Equal(t, ErrShardNotFound, err)
			},
		},
		{
			name: "get shard",
			shards: []*Shard{
				{ID: "1", Host: "localhost", Port: 8080},
				{ID: "2", Host: "localhost", Port: 8081},
			},
			ops: func(s *naive, hash hashFn) {
				shard := s.GetShard("key")
				idx := hash("key") % uint32(len(s.shards))
				require.Equal(t, s.keys[idx], shard.ID)
			},
		},
		{
			name: "get shards",
			shards: []*Shard{
				{ID: "1", Host: "localhost", Port: 8080},
				{ID: "2", Host: "localhost:", Port: 8081},
			},
			ops: func(s *naive, _ hashFn) {
				shards := make([]*Shard, 0, len(s.keys))
				for _, key := range s.keys {
					shards = append(shards, s.shards[key])
				}

				require.Equal(t, shards, s.GetShards())
			},
		},
		{
			name: "shard hash changes",
			shards: []*Shard{
				{Host: "localhost", Port: 8080, ID: "0"},
				{Host: "localhost", Port: 8081, ID: "1"},
			},
			ops: func(s *naive, hash hashFn) {
				prev := s.GetShard("key")
				idx := hash("key") % uint32(len(s.shards))
				require.Equal(t, s.keys[idx], prev.ID)

				require.NoError(t, s.RegisterShard(&Shard{Host: "localhost", Port: 8082, ID: "2"}))
				curr := s.GetShard("key")
				require.NotEqual(t, prev, curr)
			},
		},
	}

	ln := func(s string) uint32 { return uint32(len(s)) }
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			sharding := NewNaive(tc.shards, ln).(*naive)
			tc.ops(sharding, ln)
		})
	}
}
