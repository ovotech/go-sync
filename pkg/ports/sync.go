package ports

import "context"

// Sync interface can be used for downstream services that implement Sync in your own workflow.
type Sync interface {
	SyncWith(ctx context.Context, adapter Adapter) error // Sync the source service with this service.
}
