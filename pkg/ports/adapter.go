package ports

import "context"

// Adapter interfaces are used to allow Sync to communicate with third party services.
type Adapter interface {
	Get(ctx context.Context) ([]string, error) // Get things in a service.
	Add(context.Context, []string) error       // Add things to a service.
	Remove(context.Context, []string) error    // Remove things from a service.
}
