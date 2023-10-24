package eviction

import (
	"sync"
)

type lru struct {
	cap, size  int
	head, tail *Node
	cache      map[string]*Node
	mx         sync.RWMutex
}

func NewLRU(capacity int) Algorithm {
	head, tail := &Node{}, &Node{}
	head.next, tail.prev = tail, head

	return &lru{
		cap:   capacity,
		cache: make(map[string]*Node),
		head:  head,
		tail:  tail,
	}
}

func (l *lru) Get(key string) (string, bool) {
	l.mx.RLock()
	defer l.mx.RUnlock()

	if node, ok := l.cache[key]; ok {
		l.promote(node)
		return node.val, true
	}

	return "", false
}

func (l *lru) promote(node *Node) {
	left, right := node.prev, node.next
	if left != nil {
		left.next = right
	}

	if right != nil {
		right.prev = left
	}

	l.justPromote(node)
}

func (l *lru) justPromote(node *Node) {
	node.next = l.head.next
	l.head.next.prev = node
	node.prev = l.head
	l.head.next = node
}

func (l *lru) Put(key, val string) {
	l.mx.Lock()
	defer l.mx.Unlock()

	if node, ok := l.cache[key]; ok {
		node.val = val
		l.promote(node)
		return
	}

	node := &Node{key: key, val: val, prev: l.head, next: l.head.next}
	l.cache[key] = node
	l.size++
	l.justPromote(node)

	if l.size > l.cap {
		l.evict()
	}
}

func (l *lru) evict() {
	node := l.tail.prev
	prev := l.tail.prev.prev
	prev.next = l.tail
	l.tail.prev = prev

	delete(l.cache, node.key)
	l.size--
}

func (l *lru) Len() uint32 {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return uint32(l.size)
}
