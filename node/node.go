package node

import (
	"context"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/sharding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

// Nodes are used to quickly find a Node by Node.ID
type Nodes map[string]*Node

func (n Nodes) NodeIDs() []string {
	var keys = make([]string, 0, len(n))
	for _, node := range n {
		keys = append(keys, node.ID)
	}

	return keys
}

func (n Nodes) NodesApiStyle() []*api.Node {
	var nodes = make([]*api.Node, 0, len(n))
	for _, v := range n {
		nodes = append(nodes, &api.Node{
			Id:   v.ID,
			Host: v.Host,
			Port: uint32(v.Port),
		})
	}

	return nodes
}

type Node struct {
	ID   string `yaml:"id"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	// storing the gRPC client and the connection to the node, because
	// it's so expensive to create a new client every time we want to
	// send a request to the node.
	cc      *grpc.ClientConn
	gclient api.CacheServiceClient
}

func (n *Node) connString() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

// RefreshClient recreates the gRPC client, this is useful when the node
// is started and don't have a client yet.
func (n *Node) RefreshClient(ctx context.Context) error {

	// each node has a gRPC client, and it keeps connection to the node
	// alive, so we don't need to create a new client every time we want
	// to send a request to the node.
	//
	// also, we have an election mechanism to elect a leader node, and our
	// config is always get updated, before the keepalive timeout expires.
	if n.gclient != nil && n.cc != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		n.connString(),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             2 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return err
	}

	n.gclient = api.NewCacheServiceClient(conn)
	n.cc = conn
	return nil
}

// Close closes the node's gRPC client connection.
//
// we will close the connection to the node only when the node is removed
// from the config, and only when the node is down.
func (n *Node) Close() error {
	if n.cc != nil {
		return n.cc.Close()
	}

	return nil
}

func (n *Node) Request() api.CacheServiceClient {
	return n.gclient
}

func (n *Node) ToShard() *sharding.Shard {
	return &sharding.Shard{
		ID:   n.ID,
		Host: n.Host,
		Port: n.Port,
	}
}
