package eviction

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"sync"
	"testing"
)

func isValidOrder(order []Node, head *Node) bool {
	current := head.next
	for _, node := range order {
		if current.key != node.key || current.val != node.val {
			return false
		}

		current = current.next
	}

	return true
}

func TestLru_Get(t *testing.T) {
	testCases := []struct {
		name           string
		cap            int
		operate        func(t *testing.T, lru *lru)
		verifyInternal func(t *testing.T, lru *lru)
	}{
		{
			name: "get from empty lru",
			cap:  10,
			operate: func(t *testing.T, lru *lru) {
				_, ok := lru.Get("foo")
				require.False(t, ok)
			},
			verifyInternal: func(t *testing.T, lru *lru) {
				require.Equal(t, 0, lru.size)
				require.Equal(t, uint32(0), lru.Len())
			},
		},
		{
			name: "get from non-empty lru",
			cap:  10,
			operate: func(t *testing.T, lru *lru) {
				lru.Put("foo", "bar")
				_, ok := lru.Get("foo")
				require.True(t, ok)
			},
			verifyInternal: func(t *testing.T, lru *lru) {
				order := []Node{
					{key: "foo", val: "bar"},
				}

				require.Equal(t, 1, lru.size)
				require.True(t, isValidOrder(order, lru.head))
			},
		},
		{
			name: "putting the same key twice",
			cap:  10,
			operate: func(t *testing.T, lru *lru) {
				lru.Put("foo", "bar")
				lru.Put("foo", "baz")
				lru.Put("bar", "baz")
			},
			verifyInternal: func(t *testing.T, lru *lru) {
				order := []Node{
					{key: "bar", val: "baz"},
					{key: "foo", val: "baz"},
				}

				require.Equal(t, 2, lru.size)
				require.True(t, isValidOrder(order, lru.head))
			},
		},
		{
			name: "get and promote",
			cap:  10,
			operate: func(t *testing.T, lru *lru) {
				lru.Put("foo", "bar")
				lru.Put("bar", "baz")
				lru.Put("baz", "qux")
				lru.Get("foo")
			},
			verifyInternal: func(t *testing.T, lru *lru) {
				order := []Node{
					{key: "foo", val: "bar"},
					{key: "baz", val: "qux"},
					{key: "bar", val: "baz"},
				}

				require.Equal(t, 3, lru.size)
				require.True(t, isValidOrder(order, lru.head))
			},
		},
		{
			name: "evict",
			cap:  3,
			operate: func(t *testing.T, lru *lru) {
				lru.Put("foo", "bar")
				lru.Put("bar", "baz")
				lru.Put("baz", "qux")
				lru.Put("qux", "quux")
			},
			verifyInternal: func(t *testing.T, lru *lru) {
				order := []Node{
					{key: "qux", val: "quux"},
					{key: "baz", val: "qux"},
					{key: "bar", val: "baz"},
				}

				require.Equal(t, 3, lru.size)
				require.True(t, isValidOrder(order, lru.head))
			},
		},
		{
			name: "single capacity",
			cap:  1,
			operate: func(t *testing.T, lru *lru) {
				lru.Put("foo", "bar")
				lru.Put("bar", "baz")
				lru.Put("baz", "qux")
			},
			verifyInternal: func(t *testing.T, lru *lru) {
				order := []Node{
					{key: "baz", val: "qux"},
				}

				require.Equal(t, 1, lru.size)
				require.True(t, isValidOrder(order, lru.head))
			},
		},
		{
			name: "in goroutines",
			cap:  100,
			operate: func(t *testing.T, lru *lru) {
				var wg sync.WaitGroup
				for i := 0; i < lru.cap; i++ {
					wg.Add(1)

					go func(i int) {
						defer wg.Done()

						lru.Put(strconv.Itoa(i), strconv.Itoa(i))
					}(i)
				}

				wg.Wait()
			},
			verifyInternal: func(t *testing.T, lru *lru) {
				require.Equal(t, lru.cap, lru.size)
				var (
					ans = make([]string, 0, lru.cap)
					wg  sync.WaitGroup
					ch  = make(chan string, lru.cap/5)
				)

				for i := 0; i < lru.cap; i++ {
					wg.Add(1)

					go func(i int) {
						defer wg.Done()

						val, ok := lru.Get(strconv.Itoa(i))
						require.True(t, ok)
						ch <- val
					}(i)
				}

				go func() {
					defer close(ch)
					wg.Wait()
				}()

				for val := range ch {
					ans = append(ans, val)
				}

				require.Equal(t, lru.cap, len(ans))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			algo := NewLRU(tc.cap)

			ep, _ := algo.(*lru)
			tc.operate(t, ep)
			tc.verifyInternal(t, ep)
		})
	}
}
