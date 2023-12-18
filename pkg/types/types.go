package types

import "context"

// Adapter interfaces are used to allow Sync to communicate with third party services.
type Adapter interface {
	Get(ctx context.Context) (things []string, err error) // Get things in a service.
	Add(ctx context.Context, things []string) error       // Add things to a service.
	Remove(ctx context.Context, things []string) error    // Remove things from a service.
}

// ConfigKey is a configuration key to Init a new adapter.
type ConfigKey = string

// A ConfigFn is used to pass additional or custom functionality to an adapter.
type ConfigFn[T Adapter] func(T)

// InitFn is an optional adapter function that can initialise a new adapter using a static configuration.
// This is to make it easier to use an adapter in a CLI or other service that invokes adapters programmatically.
type InitFn[T Adapter] func(context.Context, map[ConfigKey]string, ...ConfigFn[T]) (T, error)
