package client

import (
	"context"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/sharding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"gopkg.in/yaml.v3"
	"hash/crc32"
	"os"
	"time"
)

// NodesConfig is used as an initial configuration for the client,
// in reality it can be changed overtime, and can be stored in some
// service discovery system like etcd,
type NodesConfig struct {

	// todo: may be use another structure
	//  shard is registered in the context of sharding algorithm
	//  here we can use another struct to make code less coupled
	Shards []*sharding.Shard `yaml:"shards"`
}

func parseConfig(path string) (*NodesConfig, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c NodesConfig
	if err = yaml.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

type client struct {
	algo sharding.Sharding
}

func NewClient(configPath string) (Client, error) {
	nodesConfig, err := parseConfig(configPath)
	if err != nil {
		return nil, err
	}

	algo := sharding.NewExpensive(
		nodesConfig.Shards,
		func(k string) uint32 { return crc32.ChecksumIEEE([]byte(k)) },
	)

	return &client{algo: algo}, nil
}

// todo: create pool of grpc clients instead of creating a new one every time
func getGrpcClient(shard *sharding.Shard) (api.CacheServiceClient, error) {
	conn, err := grpc.Dial(
		shard.ConnStr(),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		// todo: configure this properly
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    0,
			Timeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}

	return api.NewCacheServiceClient(conn), nil
}

func (c *client) Get(key string) (string, error) {
	gcl, err := getGrpcClient(c.algo.GetShard(key))
	if err != nil {
		return "", fmt.Errorf("failed to initialize grpc client: %w", err)
	}

	// todo: configure this properly
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := gcl.Get(ctx, &api.GetRequest{Key: key})
	if err != nil {
		return "", asClientError(err)
	}

	return resp.Value, nil
}

func (c *client) Put(key, value string) error {
	gcl, err := getGrpcClient(c.algo.GetShard(key))
	if err != nil {
		return fmt.Errorf("failed to initialize grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = gcl.Put(ctx, &api.PutRequest{Key: key, Value: value})
	if err != nil {
		return fmt.Errorf("failed to put value in cache: %w", err)
	}

	return nil
}
