// This package defines a LRU cache server which supports client-side consistent hashing,
// TLS (and mTLS), client access via both HTTP/gRPC,
package server

import (
	"context"
	"fmt"
	api "github.com/fadyat/speedy/api"
	"github.com/fadyat/speedy/eviction"
	"github.com/fadyat/speedy/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"go.uber.org/zap"

	"sync"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS        = "OK"
	KeyNotFoundMsg = "key not found"
)

type CacheServer struct {
	api.UnimplementedCacheServiceServer

	configPath     string
	cache          eviction.Algorithm
	router         *gin.Engine
	logger         *zap.SugaredLogger
	nodesConfig    node.NodesConfig
	leaderID       string
	nodeID         string
	groupID        string
	clientAuth     bool
	shutdownChan   chan bool
	decisionChan   chan string
	electionLock   sync.RWMutex
	electionStatus bool
	//api.UnimplementedCacheServiceServer
}

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ServerComponents struct {
	GrpcServer *grpc.Server
}

const (
	DYNAMIC = "DYNAMIC"
)

// Utility function for creating a new gRPC server secured with mTLS, and registering a cache server service with it.
// Set node_id param to DYNAMIC to dynamically discover node id.
// Otherwise, manually set it to a valid nodeID from the config file.
// Returns tuple of (gRPC server instance, registered Cache CacheServer instance).
func NewCacheServer(capacity int, configFile string, verbose bool, nodeID string, clientAuth bool) (*grpc.Server, *CacheServer) {
	// get nodes config
	nodesConfig := node.LoadNodesConfig(configFile)

	// set up logging
	sugared_logger := GetSugaredZapLogger(
		nodesConfig.ServerLogfile,
		nodesConfig.ServerErrfile,
		verbose,
	)

	// determine which node id we are and which group we are in
	var finalNodeID string
	if nodeID == DYNAMIC {
		log.Printf("passed node id: %s", nodeID)
		finalNodeID = node.GetCurrentNodeId(nodesConfig)
		log.Printf("final node id: %s", finalNodeID)

		// if this is not one of the initial nodes in the config file, add it dynamically
		if _, ok := nodesConfig.Nodes[finalNodeID]; !ok {
			host, _ := os.Hostname()
			nodesConfig.Nodes[finalNodeID] = node.NewNode(finalNodeID, host, 8080)
		}
	} else {
		finalNodeID = nodeID

		// if this is not one of the initial nodes in the config file, panic
		if _, ok := nodesConfig.Nodes[finalNodeID]; !ok {
			panic("given node ID not found in config file")
		}
	}

	// set up gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// create server instance
	cacheServer := CacheServer{
		router:       router,
		logger:       sugared_logger,
		nodesConfig:  nodesConfig,
		nodeID:       finalNodeID,
		leaderID:     NO_LEADER,
		clientAuth:   clientAuth,
		decisionChan: make(chan string, 1),
	}

	//cacheServer.router.GET("/get/:key", cacheServer.GetHandler)
	//cacheServer.router.POST("/put", cacheServer.PutHandler)
	var grpcServer *grpc.Server

	grpcServer = grpc.NewServer()

	// set up TLS
	//api.RegisterCacheServiceServer(grpcServer, &cacheServer)
	reflection.Register(grpcServer)
	return grpcServer, &cacheServer
}

func (s *CacheServer) Get(_ context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	if val, ok := s.cache.Get(req.Key); ok {
		return &api.GetResponse{Value: val}, nil
	}

	return nil, status.Error(codes.NotFound, KeyNotFoundMsg)
}

func (s *CacheServer) Put(_ context.Context, req *api.PutRequest) (*emptypb.Empty, error) {
	s.cache.Put(req.Key, req.Value)
	return &emptypb.Empty{}, nil
}

func (s *CacheServer) Len(_ context.Context, _ *emptypb.Empty) (*api.LengthResponse, error) {
	return &api.LengthResponse{Length: s.cache.Len()}, nil
}

func (s *CacheServer) GetClusterConfig(_ context.Context, _ *emptypb.Empty) (*api.ClusterConfig, error) {
	locallyStored, err := getLocallyStoredClusterConfig(s.configPath)
	if err != nil {
		return nil, err
	}

	return &api.ClusterConfig{Nodes: locallyStored}, nil
}

func getLocallyStoredClusterConfig(path string) ([]*api.Node, error) {
	// todo: can use cache here, to avoid reading from file system every time.
	//  and read only when the file is updated + update the cache.

	inode, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read initial state file: %w", err)
	}

	// todo: rewrite this peace of shit
	var nodes = struct {
		Nodes map[string]*node.Node `yaml:"nodes"`
	}{}

	if err = yaml.Unmarshal(inode, &nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial state file: %w", err)
	}

	var apiStyle = make([]*api.Node, 0, len(nodes.Nodes))
	for _, v := range nodes.Nodes {
		apiStyle = append(apiStyle, &api.Node{
			Id:   v.Id,
			Host: v.Host,
			Port: v.Port,
		})
	}

	return apiStyle, nil
}

// Utility function to get a new Cache Client which uses gRPC secured with mTLS
func (s *CacheServer) NewCacheClient(serverHost string, serverPort int) (api.CacheServiceClient, error) {

	var conn *grpc.ClientConn
	var err error

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping back
		PermitWithoutStream: true,             // send pings even without active streams
	}

	// set up connection
	addr := fmt.Sprintf("%s:%d", serverHost, serverPort)
	if err != nil {
		return nil, err
	}

	conn, err = grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithTimeout(time.Duration(time.Second)),
	)

	// set up client
	return api.NewCacheServiceClient(conn), nil
}

// Register node with the cluster. This is a function to be called internally by server code (as
// opposed to the gRPC handler to register node, which is what receives the RPC sent by this function).
func (s *CacheServer) RegisterNodeInternal() {
	s.logger.Infof("attempting to register %s with cluster", s.nodeID)
	localNode, _ := s.nodesConfig.Nodes[s.nodeID]

	// try to register with each node until one returns a successful response
	for _, node := range s.nodesConfig.Nodes {
		// skip self
		if node.Id == s.nodeID {
			continue
		}
		req := api.Node{
			Id:   localNode.Id,
			Host: localNode.Host,
			Port: localNode.Port,
		}

		c, err := s.NewCacheClient(node.Host, int(node.Port))
		if err != nil {
			s.logger.Errorf("unable to connect to node %s", node.Id)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = c.RegisterNodeWithCluster(ctx, &req)
		if err != nil {
			s.logger.Infof("error registering node %s with cluster: %v", s.nodeID, err)
			continue
		}

		s.logger.Infof("node %s is registered with cluster", s.nodeID)
		return
	}
}

// Log function that can be called externally
func (s *CacheServer) LogInfoLevel(msg string) {
	s.logger.Info(msg)
}

// Set up logger at the specified verbosity level
func GetSugaredZapLogger(logFile string, errFile string, verbose bool) *zap.SugaredLogger {
	var level zap.AtomicLevel
	outputPaths := []string{logFile}
	errorPaths := []string{errFile}

	// also log to console in verbose mode
	if verbose {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
		outputPaths = append(outputPaths, "stdout")
	} else {
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		errorPaths = append(outputPaths, "stderr")
	}
	cfg := zap.Config{
		Level:            level,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorPaths,
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
