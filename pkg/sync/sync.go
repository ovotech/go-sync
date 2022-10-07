/*
Package sync is the logic that synchronises adapters together, and determines what should be where.
*/
package sync

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ovotech/go-sync/internal/types"
	"github.com/ovotech/go-sync/pkg/ports"
)

// generateHashMap takes a list of strings and returns a hashed map of { item => true }.
func generateHashMap(i []string) map[string]bool {
	out := map[string]bool{}
	for _, str := range i {
		out[str] = true
	}

	return out
}

type Sync struct {
	DryRun bool            // DryRun mode calculates membership, but doesn't add or remove.
	Add    bool            // Perform Add operations.
	Remove bool            // Perform Remove operations.
	source ports.Adapter   // The source adapter.
	cache  map[string]bool // cache prevents polling the source more than once.
	logger types.Logger
}

// New creates a new Sync service.
func New(source ports.Adapter, optsFn ...func(*Sync)) *Sync {
	sync := &Sync{
		DryRun: false,
		Add:    true,
		Remove: true,
		source: source,
		cache:  make(map[string]bool),
		logger: log.New(os.Stderr, "[go-sync/sync] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(sync)
	}

	return sync
}

// getThingsToAdd determines things that should be added to the destination service.
func (s *Sync) getThingsToAdd(things []string) []string {
	out := make([]string, 0, len(things))
	hashMap := generateHashMap(things)

	for thing := range s.cache {
		if !hashMap[thing] {
			out = append(out, thing)
		}
	}

	return out
}

// getThingsToRemove determines things that should be removed from the destination service.
func (s *Sync) getThingsToRemove(things []string) []string {
	var out []string

	hashMap := generateHashMap(things)
	for thing := range hashMap {
		if !s.cache[thing] {
			out = append(out, thing)
		}
	}

	return out
}

// generateCache populates the cache with a map of things for efficient lookup.
func (s *Sync) generateCache(ctx context.Context) error {
	if len(s.cache) == 0 {
		s.logger.Println("Getting things from source adapter")

		things, err := s.source.Get(ctx)
		if err != nil {
			return fmt.Errorf("get -> %w", err)
		}

		s.cache = generateHashMap(things)
	}

	return nil
}

// WithLogger can be used to set a custom logger.
func WithLogger(logger types.Logger) func(*Sync) {
	return func(sync *Sync) {
		sync.logger = logger
	}
}

// remove processes removing things from a destination service.
func (s *Sync) remove(ctx context.Context, things []string, removeFn func(context.Context, []string) error) error {
	thingsToRemove := s.getThingsToRemove(things)

	if s.DryRun || !s.Remove {
		s.logger.Printf("Would remove %s, but remove is disabled", thingsToRemove)

		return nil
	}

	if len(thingsToRemove) == 0 {
		return nil
	}

	s.logger.Printf("Removing %s", thingsToRemove)

	return removeFn(ctx, thingsToRemove)
}

// add processes adding things to a destination service.
func (s *Sync) add(ctx context.Context, things []string, removeFn func(context.Context, []string) error) error {
	thingsToAdd := s.getThingsToAdd(things)

	if s.DryRun || !s.Add {
		s.logger.Printf("Would add %s, but add is disabled", thingsToAdd)

		return nil
	}

	if len(thingsToAdd) == 0 {
		return nil
	}

	s.logger.Printf("Adding %s", thingsToAdd)

	return removeFn(ctx, thingsToAdd)
}

// SyncWith synchronises the destination service with the source service, adding & removing things as necessary.
func (s *Sync) SyncWith(ctx context.Context, adapter ports.Adapter) error {
	s.logger.Println("Starting sync")

	// Call to populate the cache from the source adapter.
	if err := s.generateCache(ctx); err != nil {
		return fmt.Errorf("sync.syncwith.generateCache -> %w", err)
	}

	s.logger.Println("Getting things from destination adapter")

	things, err := adapter.Get(ctx)
	if err != nil {
		return fmt.Errorf("sync.syncwith.get -> %w", err)
	}

	s.logger.Println("Processing things to remove")

	err = s.remove(ctx, things, adapter.Remove)
	if err != nil {
		return fmt.Errorf("sync.syncwith.remove -> %w", err)
	}

	s.logger.Println("Processing things to add")

	err = s.add(ctx, things, adapter.Add)
	if err != nil {
		return fmt.Errorf("sync.syncwith.add -> %w", err)
	}

	s.logger.Println("Finished sync")

	return nil
}
