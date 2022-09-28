// Package sync provides a way to synchronise membership of arbitrary services.
package sync

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ovotech/go-sync/internal/types"
	"github.com/ovotech/go-sync/pkg/ports"
)

type Sync struct {
	DryRun bool            // Flag to indicate if running in dryRun mode.
	Add    bool            // Flag to indicate if Sync should add accounts.
	Remove bool            // Flag to indicate if Sync should remove accounts.
	source ports.Adapter   // The source of truth whose membership will be synced with other services.
	cache  map[string]bool // cache prevents polling the source more than once.
	logger types.Logger    // Custom logger.
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
		s.logger.Println("Getting accounts from source adapter")

		accounts, err := s.source.Get(ctx)
		if err != nil {
			return fmt.Errorf("get -> %w", err)
		}

		s.cache = generateHashMap(accounts)
	}

	return nil
}

// OptionLogger can be used to set a custom logger.
func OptionLogger(logger types.Logger) func(*Sync) {
	return func(sync *Sync) {
		sync.logger = logger
	}
}

// New creates a new Sync service.
func New(source ports.Adapter, optsFn ...func(*Sync)) *Sync {
	sync := &Sync{
		DryRun: false,
		Add:    true,
		Remove: true,
		source: source,
		cache:  map[string]bool{},
		logger: log.New(os.Stderr, "[go-sync/sync] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(sync)
	}

	return sync
}

// remove processes removing accounts from a destination service.
func (s *Sync) remove(ctx context.Context, accounts []string, removeFn func(context.Context, []string) error) error {
	accountsToAction := s.getAccountsToRemove(accounts)

	if s.DryRun || !s.Remove {
		s.logger.Printf("Would remove %s, but remove is disabled", accountsToAction)

		return nil
	}

	if len(accountsToAction) == 0 {
		return nil
	}

	s.logger.Printf("Removing %s", accountsToAction)

	return removeFn(ctx, accountsToAction)
}

// add processes adding accounts to a destination service.
func (s *Sync) add(ctx context.Context, accounts []string, removeFn func(context.Context, []string) error) error {
	accountsToAction := s.getAccountsToAdd(accounts)

	if s.DryRun || !s.Add {
		s.logger.Printf("Would add %s, but add is disabled", accountsToAction)

		return nil
	}

	if len(accountsToAction) == 0 {
		return nil
	}

	s.logger.Printf("Adding %s", accountsToAction)

	return removeFn(ctx, accountsToAction)
}

// SyncWith synchronises the requested service with the source service, adding & removing members.
func (s *Sync) SyncWith(ctx context.Context, adapter ports.Adapter) error {
	s.logger.Println("Starting sync")

	// Call to populate the cache from the source adapter.
	if err := s.generateCache(ctx); err != nil {
		return fmt.Errorf("sync.syncwith.generateCache -> %w", err)
	}

	s.logger.Println("Getting accounts from destination adapter")

	accounts, err := adapter.Get(ctx)
	if err != nil {
		return fmt.Errorf("sync.syncwith.get -> %w", err)
	}

	s.logger.Println("Processing accounts to remove")

	err = s.remove(ctx, accounts, adapter.Remove)
	if err != nil {
		return fmt.Errorf("sync.syncwith.remove -> %w", err)
	}

	s.logger.Println("Processing accounts to add")

	err = s.add(ctx, accounts, adapter.Add)
	if err != nil {
		return fmt.Errorf("sync.syncwith.add -> %w", err)
	}

	s.logger.Println("Finished sync")

	return nil
}
