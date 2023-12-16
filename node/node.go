// Package node contains data structures useful for representing information about nodes in a distributed system,
// as well as some utility functions for getting information about them.
package node

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	api "github.com/fadyat/speedy/api"
	"gopkg.in/yaml.v3"
	"hash/crc32"
	"log"
	"os"
	"path/filepath"
)

// NodesConfig struct holds info about all server nodes in the network
type NodesConfig struct {
	Nodes            map[string]*Node `json:"nodes"`
	EnableClientAuth bool             `json:"enable_client_auth"`
	EnableHttps      bool             `json:"enable_https"`
	ServerLogfile    string           `json:"server_logfile"`
	ServerErrfile    string           `json:"server_errfile"`
	ClientLogfile    string           `json:"client_logfile"`
	ClientErrfile    string           `json:"client_errfile"`
}

// Node struct contains all info we need about a server node, as well as
// a gRPC client to interact with it
type Node struct {
	Id         string `json:"id"`
	Host       string `json:"host"`
	Port       uint32 `json:"port"`
	HashId     uint32
	GrpcClient api.CacheServiceClient
}

func (n *Node) SetGrpcClient(c api.CacheServiceClient) {
	n.GrpcClient = c
}

func NewNode(id, host string, port uint32) *Node {
	return &Node{
		Id:     id,
		Host:   host,
		Port:   port,
		HashId: HashId(id),
	}
}

type Nodes []*Node

// Len implementing methods required for sorting
func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].HashId < n[j].HashId }

func HashId(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// LoadNodesConfig Load nodes config file
func LoadNodesConfig(configFile string) NodesConfig {
	file, _ := os.ReadFile(filepath.Clean(configFile))

	// set defaults
	nodesConfig := NodesConfig{
		EnableClientAuth: true,
		EnableHttps:      true,
		ServerLogfile:    "minicache.log",
		ServerErrfile:    "minicache.err",
		ClientLogfile:    "minicache-client.log",
		ClientErrfile:    "minicache-client.err",
	}

	_ = yaml.Unmarshal(file, &nodesConfig)

	// if config is empty, add 1 node at localhost:8080
	if len(nodesConfig.Nodes) == 0 {
		log.Printf("couldn't find config file or it was empty: %s", configFile)
		log.Println("using default node localhost")
		nodesConfig = NodesConfig{Nodes: make(map[string]*Node)}
		defaultNode := NewNode("node0", "localhost", 8080)
		nodesConfig.Nodes[defaultNode.Id] = defaultNode
	} else {
		for _, nodeInfo := range nodesConfig.Nodes {
			nodeInfo.HashId = HashId(nodeInfo.Id)
		}
	}
	return nodesConfig
}

// GetCurrentNodeId Determine which node ID we are
func GetCurrentNodeId(config *NodesConfig) string {
	host, _ := os.Hostname()
	for _, node := range config.Nodes {
		if node.Host == host {
			return node.Id
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
