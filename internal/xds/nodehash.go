package xds

import core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

type IDHash struct{}

func (IDHash) ID(node *core.Node) string {
	return "velo-global"
}
