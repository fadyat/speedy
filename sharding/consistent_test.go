package sharding

import (
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

func TestConsistent_Flow(t *testing.T) {
	testcases := []struct {
		name   string
		shards []*Shard
		ops    func(s *consistent, fn hashFn)
	}{
		{
			name: "new consistent",
			shards: []*Shard{
				{ID: "10"},
				{ID: "200"},
				{ID: "3"},
				{ID: "3"},
			},
			ops: func(s *consistent, _ hashFn) {
				require.Equal(t, 3, len(s.shards))
				require.Equal(t, 3, len(s.orderedKeys))
				require.Equal(t, uint32(1), s.orderedKeys[0])
				require.Equal(t, uint32(2), s.orderedKeys[1])
				require.Equal(t, uint32(3), s.orderedKeys[2])
				require.Equal(t, s.shards[s.orderedKeys[0]], &Shard{ID: "3"})
			},
		},
		{
			name: "register shard",
			ops: func(s *consistent, _ hashFn) {
				require.NoError(t, s.RegisterShard(&Shard{ID: "1"}))

				require.Equal(t, 1, len(s.shards))
				require.Equal(t, 1, len(s.orderedKeys))
				require.Equal(t, uint32(1), s.orderedKeys[0])
			},
		},
		{
			name: "register multiple shards",
			ops: func(s *consistent, _ hashFn) {
				require.NoError(t, s.RegisterShard(&Shard{ID: "1000"}))
				require.NoError(t, s.RegisterShard(&Shard{ID: "2"}))

				require.Equal(t, 2, len(s.shards))
				require.Equal(t, 2, len(s.orderedKeys))
				require.Equal(t, uint32(1), s.orderedKeys[0])
				require.Equal(t, uint32(4), s.orderedKeys[1])
			},
		},
		{
			name:   "register shard twice",
			shards: []*Shard{{ID: "1"}},
			ops: func(s *consistent, _ hashFn) {
				err := s.RegisterShard(&Shard{ID: "1"})
				require.Equal(t, ErrShardAlreadyRegistered, err)
			},
		},
		{
			name: "delete shard",
			shards: []*Shard{
				{ID: "10"},
				{ID: "200"},
				{ID: "3"},
			},
			ops: func(s *consistent, _ hashFn) {
				require.NoError(t, s.DeleteShard(&Shard{ID: "10"}))
				require.Equal(t, 2, len(s.shards))
				require.Equal(t, 2, len(s.orderedKeys))
				require.Equal(t, uint32(1), s.orderedKeys[0])
				require.Equal(t, uint32(3), s.orderedKeys[1])
			},
		},
		{
			name: "delete shard not found",
			ops: func(s *consistent, _ hashFn) {
				err := s.DeleteShard(&Shard{ID: "10"})
				require.Equal(t, ErrShardNotFound, err)
			},
		},
		{
			name: "get shard",
			shards: []*Shard{
				{ID: "1000000000"},
				{ID: "200"},
				{ID: "3"},
			},
			ops: func(s *consistent, hash hashFn) {
				// exact match
				shard := s.GetShard("key")
				require.Equal(t, &Shard{ID: "200"}, shard)

				// closest match, clockwise
				shard = s.GetShard("key1")
				require.Equal(t, &Shard{ID: "1000000000"}, shard)

				// overflow match
				shard = s.GetShard(">1000000000")
				require.Equal(t, &Shard{ID: "3"}, shard)
			},
		},
		{
			name: "get shards",
			shards: []*Shard{
				{ID: "1000000000"},
				{ID: "200"},
			},
			ops: func(s *consistent, _ hashFn) {
				shards := s.GetShards()
				for i := range s.orderedKeys {
					storedShard := s.shards[s.orderedKeys[i]]
					require.Equal(t, storedShard, shards[i])
				}
			},
		},
		{
			name: "not reassigned when adding shard",
			shards: []*Shard{
				{ID: "bebrabebra"},
				{ID: "200"},
				{ID: "3"},
			},
			ops: func(s *consistent, hash hashFn) {
				prevShardStable := s.GetShard("ab")
				prevShardMutable := s.GetShard("aboba")
				require.Equal(t, &Shard{ID: "200"}, prevShardStable)
				require.Equal(t, &Shard{ID: "bebrabebra"}, prevShardMutable)

				newShard := &Shard{ID: "aboba"}
				require.NoError(t, s.RegisterShard(newShard))
				require.Equal(t, prevShardStable, s.GetShard("ab"))
				require.NotEqual(t, prevShardMutable, s.GetShard("aboba"))
				require.Equal(t, newShard, s.GetShard("aboba"))
			},
		},
	}

	ln := func(key string) uint32 {
		key = strings.TrimPrefix(key, "node")
		return uint32(len(key))
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewConsistent(tc.shards, ln).(*consistent)
			tc.ops(s, ln)

			require.True(t, sort.SliceIsSorted(s.orderedKeys, func(i, j int) bool {
				return s.orderedKeys[i] < s.orderedKeys[j]
			}))
		})
	}
}
