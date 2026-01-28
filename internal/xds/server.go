package xds

import (
	"context"
	"fmt"
	"log"
	"net"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
)

const (
	defaultPort     = 18000
	defaultNodeHash = "velo-global"
	defaultVersion  = "1"
)

type ControlPlane struct {
	cache      cache.SnapshotCache
	xdsServer  server.Server
	grpcServer *grpc.Server
	logger     *log.Logger
}

func NewControlPlane(ctx context.Context, baseLogger *log.Logger) (*ControlPlane, error) {
	if baseLogger == nil {
		baseLogger = log.Default()
	}

	cacheLogger := NewLogger(baseLogger)
	snapshotCache := cache.NewSnapshotCache(true, IDHash{}, cacheLogger)
	callbacks := NewCallbacks(baseLogger)
	xdsServer := server.NewServer(ctx, snapshotCache, callbacks)

	resources := map[resource.Type][]types.Resource{
		resource.ClusterType:  {},
		resource.EndpointType: {},
		resource.ListenerType: {},
		resource.RouteType:    {},
	}
	snapshot, err := cache.NewSnapshot(defaultVersion, resources)
	if err != nil {
		return nil, err
	}

	if err := snapshotCache.SetSnapshot(ctx, defaultNodeHash, snapshot); err != nil {
		return nil, err
	}

	return &ControlPlane{
		cache:     snapshotCache,
		xdsServer: xdsServer,
		logger:    baseLogger,
	}, nil
}

func (c *ControlPlane) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		return err
	}

	c.grpcServer = grpc.NewServer()
	registerXDSServices(c.grpcServer, c.xdsServer)

	c.logger.Printf("xds control plane listening on %d", defaultPort)
	return c.grpcServer.Serve(lis)
}

func (c *ControlPlane) Stop() {
	if c.grpcServer != nil {
		c.grpcServer.GracefulStop()
	}
}

func registerXDSServices(grpcServer *grpc.Server, srv server.Server) {
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, srv)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, srv)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, srv)
}
