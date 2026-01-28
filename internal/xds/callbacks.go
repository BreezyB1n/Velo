package xds

import (
	"context"
	"log"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type Callbacks struct {
	logger *log.Logger
}

func NewCallbacks(logger *log.Logger) *Callbacks {
	return &Callbacks{logger: logger}
}

func (c *Callbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	c.logger.Printf("xds stream opened: id=%d type=%s", streamID, typeURL)
	return nil
}

func (c *Callbacks) OnStreamClosed(streamID int64, node *core.Node) {
	nodeID := nodeIDFromNode(node)
	c.logger.Printf("xds stream closed: id=%d node_id=%s", streamID, nodeID)
}

func (c *Callbacks) OnStreamRequest(streamID int64, req *discovery.DiscoveryRequest) error {
	nodeID := nodeIDFromRequest(req)
	c.logger.Printf(
		"xds stream request: id=%d node_id=%s type=%s version=%s resources=%d",
		streamID,
		nodeID,
		req.GetTypeUrl(),
		req.GetVersionInfo(),
		len(req.GetResourceNames()),
	)
	return nil
}

func (c *Callbacks) OnStreamResponse(ctx context.Context, streamID int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	nodeID := nodeIDFromRequest(req)
	c.logger.Printf(
		"xds stream response: id=%d node_id=%s type=%s version=%s resources=%d",
		streamID,
		nodeID,
		req.GetTypeUrl(),
		resp.GetVersionInfo(),
		len(resp.GetResources()),
	)
}

func (c *Callbacks) OnDeltaStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	return nil
}

func (c *Callbacks) OnDeltaStreamClosed(streamID int64, node *core.Node) {
}

func (c *Callbacks) OnStreamDeltaRequest(streamID int64, req *discovery.DeltaDiscoveryRequest) error {
	return nil
}

func (c *Callbacks) OnStreamDeltaResponse(streamID int64, req *discovery.DeltaDiscoveryRequest, resp *discovery.DeltaDiscoveryResponse) {
}

func (c *Callbacks) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	return nil
}

func (c *Callbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
}

func nodeIDFromNode(node *core.Node) string {
	if node == nil {
		return ""
	}
	return node.GetId()
}

func nodeIDFromRequest(req *discovery.DiscoveryRequest) string {
	if req == nil {
		return ""
	}
	return nodeIDFromNode(req.GetNode())
}
