/*
Package ports lists the types of methods expected from Sync's adapters.
*/
package ports

import "context"

// Adapter interfaces are used to allow Sync to communicate with third party services.
type Adapter interface {
	Get(ctx context.Context) ([]string, error) // Get things in a service.
	Add(context.Context, []string) error       // Add things to a service.
	Remove(context.Context, []string) error    // Remove things from a service.
}

// The Sync interface can be used for downstream services that implement Sync in your own workflow.
type Sync interface {
	SyncWith(ctx context.Context, adapter Adapter) error // Sync the things in a source service with this service.
}
