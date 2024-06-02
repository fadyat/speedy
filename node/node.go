// Package node contains data structures useful for representing information about nodes in a distributed system,
// as well as some utility functions for getting information about them.
package node

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/sharding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Node struct contains all info we need about a server node, as well as
// a gRPC client to interact with it
type Node struct {
	ID   string `yaml:"id"`
	Host string `yaml:"host"`
	Port uint32 `yaml:"port"`

	// storing the gRPC client and the connection to the node, because
	// it's so expensive to create a new client every time we want to
	// send a request to the node.
	cc      *grpc.ClientConn
	gclient api.CacheServiceClient
}

func (n *Node) SetGrpcClient(c api.CacheServiceClient) {
	n.gclient = c
}

func NewNode(id, host string, port uint32) *Node {
	return &Node{
		ID:   id,
		Host: host,
		Port: port,
	}
}

func (n *Node) connString() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

type Nodes map[string]*Node

func LoadNodesConfig(configFile string) *NodesConfig {
	file, _ := os.ReadFile(filepath.Clean(configFile))

	// set defaults
	nodesConfig := NodesConfig{
		ServerLogfile: "minicache.log",
		ServerErrfile: "minicache.err",
		ClientLogfile: "minicache-client.log",
		ClientErrfile: "minicache-client.err",
	}

	_ = yaml.Unmarshal(file, &nodesConfig)

	// if config is empty, add 1 node at localhost:8080
	if len(nodesConfig.Nodes) == 0 {
		log.Printf("couldn't find config file or it was empty: %s", configFile)
		log.Println("using default node localhost")
		nodesConfig = NodesConfig{Nodes: make(map[string]*Node)}
		defaultNode := NewNode("node0", "localhost", 8080)
		nodesConfig.Nodes[defaultNode.ID] = defaultNode
	}
	return &nodesConfig
}

// GetCurrentNodeId Determine which node ID we are
func GetCurrentNodeId(config *NodesConfig) string {
	host, _ := os.Hostname()
	for _, node := range config.Nodes {
		if node.Host == host {
			return node.ID
		}
	}
	// if host not found, generate random node id
	return randSeq(5)
}

// GetRandomNode Get random node from a given list of nodes
func GetRandomNode(nodes []*Node) *Node {
	randomIndex, err := cryptoRandInt(len(nodes))
	if err != nil {
		return nil
	}
	return nodes[randomIndex]
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
		grpc.WithTransportCredentials(insecure.NewCredentials()),

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

func randSeq(n int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := cryptoRandBytes(n)
	if err != nil {
		return ""
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return fmt.Sprintf("node-%s", string(bytes))
}

func cryptoRandInt(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("invalid argument to cryptoRandInt")
	}

	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint64(b[:]) % uint64(max)), nil
}

func cryptoRandBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid argument to cryptoRandBytes")
	}

	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

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
