package gosync

import "context"

// Adapter interfaces are used to allow Sync to communicate with third party services.
type Adapter interface {
	Get(ctx context.Context) (things []string, err error) // Get things in a service.
	Add(ctx context.Context, things []string) error       // Add things to a service.
	Remove(ctx context.Context, things []string) error    // Remove things from a service.
}

// Service can be used for downstream services that implement Sync in your own workflow.
type Service interface {
	SyncWith(ctx context.Context, adapter Adapter) error // Sync the things in a source service with this service.
}
