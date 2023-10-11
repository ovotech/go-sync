package plugin

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	gosyncplugin "github.com/ovotech/go-sync/internal/plugin"
	"github.com/ovotech/go-sync/pkg/types"
)

// Serve your Go Sync plugins.
func Serve[T types.Adapter](initFn types.InitFn[T], configFns ...types.ConfigFn[T]) {
	serveConfig := &plugin.ServeConfig{
		GRPCServer:      plugin.DefaultGRPCServer,
		HandshakeConfig: gosyncplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			gosyncplugin.AdapterName: &gosyncplugin.Plugin{
				InitFn: func(ctx context.Context, config map[types.ConfigKey]string) (types.Adapter, error) {
					return initFn(ctx, config, configFns...)
				},
			},
		},
		Logger: hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		}),
	}

	plugin.Serve(serveConfig)
}
