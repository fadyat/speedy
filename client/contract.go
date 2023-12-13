package client

import (
	"context"
)

// Client is the interface that wraps the basic methods of a cache client.
type Client interface {
	Get(key string) (string, error)
	Put(key, value string) error

	// SyncClusterConfig under the hood, periodically goes to the server and
	// fetches the latest cluster configuration.
	//
	// Current and desired configs are compared, and if there is a difference,
	// the client updates its internal state:
	// - opening new connections to new nodes
	// - closing connections to nodes that are no longer part of the cluster
	SyncClusterConfig(ctx context.Context) <-chan error
}
