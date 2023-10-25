package sharding

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRendezvous_Flow(t *testing.T) {
	testcases := []struct {
		name   string
		shards []*Shard
		hashFn func(key string) uint32
		ops    func(s *rendezvous, fn hashFn)
	}{
		{
			name: "register shard",
			ops: func(s *rendezvous, _ hashFn) {
				require.NoError(t, s.RegisterShard(&Shard{Host: "localhost", Port: 8080}))
				require.Equal(t, 1, len(s.shards))
			},
		},
		{
			name: "register shard twice",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
			},
			ops: func(s *rendezvous, _ hashFn) {
				err := s.RegisterShard(&Shard{Host: "localhost", Port: 8080})
				require.Equal(t, ErrShardAlreadyRegistered, err)
			},
		},
		{
			name: "delete shard",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
			},
			ops: func(s *rendezvous, _ hashFn) {
				require.NoError(t, s.DeleteShard(&Shard{Host: "localhost", Port: 8080}))
				require.Equal(t, 0, len(s.shards))
			},
		},
		{
			name: "delete shard not found",
			ops: func(s *rendezvous, _ hashFn) {
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
			hashFn: func(key string) uint32 {
				return 0
			},
			ops: func(s *rendezvous, hash hashFn) {
				shard := s.GetShard("key")
				require.Equal(t, s.shards[0], shard)
			},
		},
		{
			name: "get shards",
			shards: []*Shard{
				{Host: "localhost", Port: 8080},
				{Host: "localhost:", Port: 8081},
			},
			ops: func(s *rendezvous, _ hashFn) {
				shards := s.GetShards()
				require.Equal(t, s.shards, shards)
			},
		},
		{
			name: "shard still consistent after adding/removing shard",
			shards: []*Shard{
				{Host: "localhost", Port: 8080, ID: "1"},
				{Host: "localhost", Port: 8081, ID: "2"},
				{Host: "localhost", Port: 8082, ID: "3"},
			},
			hashFn: func(key string) uint32 {
				switch key {
				case "key1":
					return 1
				case "key2":
					return 2
				}

				return 0
			},
			ops: func(s *rendezvous, hash hashFn) {
				prev := s.GetShard("key")
				require.Equal(t, prev, s.shards[1])

				require.NoError(t, s.RegisterShard(&Shard{Host: "localhost", Port: 8083, ID: "4"}))
				curr := s.GetShard("key")
				require.Equal(t, prev, curr)

				require.NoError(t, s.DeleteShard(&Shard{Host: "localhost", Port: 8083, ID: "4"}))
				curr = s.GetShard("key")
				require.Equal(t, prev, curr)

				require.NoError(t, s.DeleteShard(&Shard{Host: "localhost", Port: 8082, ID: "3"}))
				curr = s.GetShard("key")
				require.Equal(t, prev, curr)
			},
		},
		{
			name: "shard hash changes after adding/removing shard",
			shards: []*Shard{
				{Host: "localhost", Port: 8080, ID: "1"},
			},
			hashFn: func(key string) uint32 {
				switch key {
				case "key2":
					return 1
				}

				return 0
			},
			ops: func(s *rendezvous, hash hashFn) {
				prev := s.GetShard("key")
				require.Equal(t, prev, s.shards[0])

				require.NoError(t, s.RegisterShard(&Shard{Host: "localhost", Port: 8081, ID: "2"}))
				curr := s.GetShard("key")
				require.NotEqual(t, prev, curr)

				require.NoError(t, s.DeleteShard(&Shard{Host: "localhost", Port: 8081, ID: "2"}))
				curr = s.GetShard("key")
				require.Equal(t, prev, curr)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			sharding := NewRendezvous(tc.shards, tc.hashFn).(*rendezvous)
			tc.ops(sharding, tc.hashFn)
		})
	}
}
