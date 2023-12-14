package syncer

import (
	"context"
	"errors"
	"fmt"
	"github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const (
	ActiveTimeout = 10 * time.Second
	Timeout       = 5 * time.Second

	PingTimeout = 1 * time.Second
)

type ServerSyncer struct {
	configPath string
	syncer     *Cluster
	masterApi  api.CacheServiceClient
	masterInfo *NodeConfig
}

func sliceToMap(nodes []*api.Node) map[string]*NodeConfig {
	var m = make(map[string]*NodeConfig)
	for _, n := range nodes {
		m[n.Id] = NewNodeConfig(n.Id, n.Host, int(n.Port))
	}

	return m
}

func NewServerSyncer(configPath string) *ServerSyncer {
	var s = &ServerSyncer{
		configPath: configPath,

		// masterApi is lazy initialized, when cluster
		// state will be stable, connection will be established.
		masterApi: nil,
	}

	var (
		currentStateFetcher = func(ctx context.Context) (*CacheConfig, error) {
			return pkg.FromYaml[CacheConfig](configPath)
		}
		desiredStateFetcher = func(ctx context.Context) (*CacheConfig, error) {
			if s.masterApi == nil {
				return nil, errors.New("master api is not initialized yet")
			}

			c, err := s.masterApi.GetClusterConfig(ctx, &emptypb.Empty{})
			if err != nil {
				return nil, err
			}

			return &CacheConfig{Nodes: sliceToMap(c.Nodes)}, nil
		}
		syncStates = func(ctx context.Context, diff *cacheConfigDiff) (bool, error) {
			return applyChangesToConfigFile(ctx, diff, configPath)
		}
	)

	s.syncer = NewCluster(currentStateFetcher, desiredStateFetcher, syncStates)
	return s
}

func (s *ServerSyncer) IsMasterReady() (bool, error) {

	// when leader election algorithm will work, all nodes will store
	// master information in the config file, so we can read it from there.
	cfg, err := pkg.FromYaml[CacheConfig](s.configPath)
	if err != nil {
		return false, fmt.Errorf("failed to read cluster config: %w", err)
	}

	// cluster not initialized yet.
	if cfg.MasterInfo == nil && s.masterApi == nil {
		return false, nil
	}

	// skipping this, because it can still work without master info,
	// if previous configuration is still valid.
	if cfg.MasterInfo == nil && s.masterApi != nil {
		return s.isMasterAlive(context.Background())
	}

	// if master info is the same, we can skip initialization.
	if cfg.MasterInfo.Same(s.masterInfo) {
		return s.isMasterAlive(context.Background())
	}

	s.masterApi, err = newMasterApi(cfg.MasterInfo.Host, cfg.MasterInfo.Port)
	if err != nil {
		return false, fmt.Errorf("failed to initialize master api: %w", err)
	}

	s.masterInfo = cfg.MasterInfo
	return true, nil
}

func (s *ServerSyncer) isMasterAlive(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, PingTimeout)
	defer cancel()

	_, err := s.masterApi.Ping(ctx, &emptypb.Empty{})
	if err == nil {
		return true, nil
	}

	// forcing reinitialization of master api
	s.masterApi = nil
	return false, err
}

func newMasterApi(host string, port int) (api.CacheServiceClient, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                ActiveTimeout,
			Timeout:             Timeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial master: %w", err)
	}

	return api.NewCacheServiceClient(conn), nil
}
