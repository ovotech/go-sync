// Package sync provides a way to synchronise membership of arbitrary services.
package sync

import (
	"context"
	"fmt"

	"github.com/ovotech/go-sync/pkg/ports"
)

type Sync struct {
	DryRun bool            // Flag to indicate if running in dryRun mode.
	Add    bool            // Flag to indicate if Sync should add accounts.
	Remove bool            // Flag to indicate if Sync should remove accounts.
	source ports.Adapter   // The source of truth whose membership will be synced with other services.
	cache  map[string]bool // cache prevents polling the source more than once.
}

// generateHashMap takes a list of strings and returns a hashed map of { item => true }.
func generateHashMap(i []string) map[string]bool {
	out := map[string]bool{}
	for _, str := range i {
		out[str] = true
	}

	return out
}

// getAccountsToAdd takes a list of accounts, and returns a list that aren't in the cache.
func (s *Sync) getAccountsToAdd(accounts []string) []string {
	out := make([]string, 0, len(accounts))
	hashMap := generateHashMap(accounts)

	for account := range s.cache {
		if !hashMap[account] {
			out = append(out, account)
		}
	}

	return out
}

// getAccountsToRemove takes a list of accounts, and returns a list that aren't in the cache.
func (s *Sync) getAccountsToRemove(accounts []string) []string {
	var out []string

	hashMap := generateHashMap(accounts)
	for account := range hashMap {
		if !s.cache[account] {
			out = append(out, account)
		}
	}

	return out
}

// generateCache populates the cache with a hash of accounts in the source for efficient lookup.
func (s *Sync) generateCache(ctx context.Context) error {
	if len(s.cache) == 0 {
		accounts, err := s.source.Get(ctx)
		if err != nil {
			return fmt.Errorf("get -> %w", err)
		}

		s.cache = generateHashMap(accounts)
	}

	return nil
}

// New creates a new Sync service.
func New(source ports.Adapter, optsFn ...func(*Sync)) *Sync {
	sync := &Sync{
		DryRun: false,
		Add:    true,
		Remove: true,
		source: source,
		cache:  map[string]bool{},
	}

	for _, fn := range optsFn {
		fn(sync)
	}

	return sync
}

// perform runs the diff func, and then actions it to return the response.
func (s *Sync) perform(
	ctx context.Context,
	accounts []string,
	diff func([]string) []string,
	action func(context.Context, []string) error,
) error {
	accountsToAction := diff(accounts)

	if s.DryRun {
		// If running in dry-run mode, return the diff as success, but don't action the change.
		return nil
	}

	// If nothing needs to happen, then just return here with empty values.
	if len(accountsToAction) == 0 {
		return nil
	}

	return action(ctx, accountsToAction)
}

// SyncWith synchronises the requested service with the source service, adding & removing members.
func (s *Sync) SyncWith(ctx context.Context, adapter ports.Adapter) error {
	// Call to populate the cache from the source adapter.
	if err := s.generateCache(ctx); err != nil {
		return fmt.Errorf("sync.syncwith.generateCache -> %w", err)
	}

	accounts, err := adapter.Get(ctx)
	if err != nil {
		return fmt.Errorf("sync.syncwith.get -> %w", err)
	}

	if s.Remove {
		err = s.perform(ctx, accounts, s.getAccountsToRemove, adapter.Remove)
		if err != nil {
			return fmt.Errorf("sync.syncwith.remove -> %w", err)
		}
	}

	if s.Add {
		err = s.perform(ctx, accounts, s.getAccountsToAdd, adapter.Add)
		if err != nil {
			return fmt.Errorf("sync.syncwith.add -> %w", err)
		}
	}

	return nil
}
