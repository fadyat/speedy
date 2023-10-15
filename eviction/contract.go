package eviction

var (
	noKey          string
	KeyNotFoundMsg = "key not found"
)

type Node struct {
	key, val   string
	prev, next *Node
}

type Algorithm interface {

	// Get returns the value associated with the given key.
	// - If the key exists, it returns the value and true, otherwise it returns
	//   nil and false.
	Get(key string) (string, bool)

	// Put inserts the given key-value pair into the cache.
	// - If the key already exists, it updates the value.
	// - If the key does not exist, it inserts the key-value pair.
	// - If the cache is full, it evicts the least recently used item before
	//   inserting the new key-value pair.
	Put(key, val string)
}
