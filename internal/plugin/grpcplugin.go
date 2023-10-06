package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/ovotech/go-sync/internal/proto"
	"github.com/ovotech/go-sync/pkg/types"
)

type UntypedInitFn = func(ctx context.Context, config map[types.ConfigKey]string) (types.Adapter, error)

// Ensure Plugin matches the GRPCPlugin spec by Hashicorp.
var _ plugin.GRPCPlugin = &Plugin{}

type Plugin struct {
	plugin.Plugin

	Name   string
	InitFn UntypedInitFn
}

func (g *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterAdapterServer(s, &Server{InitFn: g.InitFn})

	return nil
}

func (g *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &Client{AdapterClient: proto.NewAdapterClient(c)}, nil
}
