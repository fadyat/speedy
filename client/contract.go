package client

// Client is the interface that wraps the basic methods of a cache client.
type Client interface {
	Get(key string) (string, error)
	Put(key, value string) error
}
