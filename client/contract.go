package client

// Client is the interface that wraps the basic methods of a cache client.
//
// It's a little bit different from Algorithm in eviction/contract.go, because
// of network communication.
type Client interface {
	Get(key string) (string, error)
	Put(key, value string) error
}
