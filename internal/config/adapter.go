package config

import (
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	gosyncplugin "github.com/ovotech/go-sync/internal/plugin"
	"github.com/ovotech/go-sync/pkg/types"
)

type Adapter struct {
	Plugin string                     `validate:"required,file" yaml:"plugin"`
	Config map[types.ConfigKey]string `yaml:"config,omitempty"`
}

// ToAdapter returns an initialised Go Sync adapter.
func (a *Adapter) ToAdapter(ctx context.Context) (types.Adapter, func(), error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  gosyncplugin.HandshakeConfig,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Plugins: map[string]plugin.Plugin{
			gosyncplugin.AdapterName: &gosyncplugin.Plugin{},
		},
		Cmd: exec.Command(a.Plugin), //nolint:gosec

		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,

		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: os.Stdout,
			Level:  hclog.Debug,
		}),
	})

	// Return a GRPC Client for communicating with the adapter.
	grpcClient, err := client.Client()
	if err != nil {
		client.Kill()

		return nil, nil, fmt.Errorf("adapter.client -> %w", err)
	}

	// Dispense an adapter from the client using the standard adapter name.
	dispensedAdapter, err := grpcClient.Dispense(gosyncplugin.AdapterName)
	if err != nil {
		client.Kill()

		return nil, nil, fmt.Errorf("adapter.dispense -> %w", err)
	}

	// Coerce the adapter into a Go Sync plugin client.
	// This is identical to an adapter, except with an Init function that will be called below.
	adapter := dispensedAdapter.(*gosyncplugin.Client)

	// Get the environment variables that start with GOSYNC_*.
	environmentVariables := getEnvironmentVariables()

	// Copy the config over the environment variables.
	maps.Copy(environmentVariables, a.Config)

	// Initialise the plugin with the requested configuration.
	err = adapter.Init(ctx, environmentVariables)
	if err != nil {
		client.Kill()

		return nil, nil, fmt.Errorf("adapter.init -> %w", err)
	}

	return adapter, client.Kill, nil
}
