// Package ports lists the types of methods expected from Sync's standard adapters.
package ports

import "context"

// Adapter interfaces are used to allow Sync to communicate with third party services.
type Adapter interface {
	Get(ctx context.Context) ([]string, error) // Get users in a service.
	Add(context.Context, []string) error       // Add users to a service.
	Remove(context.Context, []string) error    // Remove users from a service.
}
