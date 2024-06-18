package gosync

import (
	"context"
	"fmt"
	"github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/types"
	"log"
	"os"
	"strings"
)

// OperatingMode specifies how Sync operates, which sync operations are run and in what order.
type OperatingMode string

const (
	// AddOnly only runs add operations.
	AddOnly OperatingMode = "Add"
	// RemoveOnly only runs remove operations.
	RemoveOnly OperatingMode = "Remove"
	// RemoveAdd first removes things, then adds them.
	RemoveAdd OperatingMode = "RemoveAdd"
	// AddRemove first adds things, then removes them.
	AddRemove OperatingMode = "AddRemove"
	// NoChangeLimit tells Sync not to set a change limit.
	NoChangeLimit int = -1
)

type Sync struct {
	DryRun        bool            // DryRun mode calculates membership, but doesn't add or remove.
	OperatingMode OperatingMode   // Change the order of Sync's operation. Default is RemoveAdd.
	CaseSensitive bool            // CaseSensitive sets if Go Sync is case-sensitive. Default is true.
	source        types.Adapter   // The source adapter.
	cache         map[string]bool // cache prevents polling the source more than once.
	/*
		MaximumChanges sets the maximum number of allowed changes per add/remove operation. It is not a cumulative
		total, and the number only applies to each distinct operation.

		For example:

		Setting this value to 3 means that a maximum of 3 things can be added AND removed from a destination (total 6)
		changes before Sync returns an ErrTooManyChanges error.

		Default is NoChangeLimit (or -1).
	*/
	MaximumChanges int
	Logger         *log.Logger
}

// New creates a new Sync service.
func New(source types.Adapter, optsFn ...func(*Sync)) *Sync {
	sync := &Sync{
		DryRun:         false,
		OperatingMode:  RemoveAdd,
		CaseSensitive:  true,
		source:         source,
		cache:          make(map[string]bool),
		MaximumChanges: NoChangeLimit,
		Logger:         log.New(os.Stderr, "[go-sync/sync] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(sync)
	}

	return sync
}

// generateHashMap takes a list of strings and returns a hashed map of { item => true }.
func (s *Sync) generateHashMap(i []string) map[string]bool {
	out := map[string]bool{}

	for _, str := range i {
		if s.CaseSensitive {
			out[str] = true
		} else {
			out[strings.ToLower(str)] = true
		}
	}

	return out
}

// getThingsToAdd determines things that should be added to the destination service.
func (s *Sync) getThingsToAdd(things []string) []string {
	out := make([]string, 0, len(things))
	hashMap := s.generateHashMap(things)

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

	hashMap := s.generateHashMap(things)
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
		s.Logger.Println("Getting things from source adapter")

		things, err := s.source.Get(ctx)
		if err != nil {
			return fmt.Errorf("get -> %w", err)
		}

		s.cache = s.generateHashMap(things)
	}

	return nil
}

// perform processes adding/removing things from a destination service.
func (s *Sync) perform(
	ctx context.Context,
	action string,
	things []string,
	diffFn func(things []string) []string,
	executeFn func(context.Context, []string) error,
) func() error {
	return func() error {
		s.Logger.Printf("Processing things to %s\n", action)

		thingsToChange := diffFn(things)

		// If the changes exceed the maximum change limit, fail with the ErrTooManyChanges error.
		if len(thingsToChange) > s.MaximumChanges && s.MaximumChanges != NoChangeLimit {
			return fmt.Errorf("%s(%v) -> %w(%v)", action, thingsToChange, errors.ErrTooManyChanges, s.MaximumChanges)
		}

		if s.DryRun {
			s.Logger.Printf("Would %s %s, but running in dry run mode", action, thingsToChange)

			return nil
		}

		if len(thingsToChange) == 0 {
			return nil
		}

		s.Logger.Printf("%s: %s", action, thingsToChange)

		err := executeFn(ctx, thingsToChange)
		if err != nil {
			return fmt.Errorf("%s(%v) -> %w", action, things, err)
		}

		return nil
	}
}

// SyncWith synchronises the destination service with the source service, adding & removing things as necessary.
func (s *Sync) SyncWith(ctx context.Context, adapter types.Adapter) error {
	s.Logger.Println("Starting sync")

	// Call to populate the cache from the source adapter.
	if err := s.generateCache(ctx); err != nil {
		return fmt.Errorf("sync.syncwith.generateCache -> %w", err)
	}

	s.Logger.Println("Getting things from destination adapter")

	things, err := adapter.Get(ctx)
	if err != nil {
		return fmt.Errorf("sync.syncwith.get -> %w", err)
	}

	s.Logger.Printf("Running in %s operating mode", s.OperatingMode)

	operations := make([]func() error, 0, 2) //nolint:gomnd,mnd

	switch s.OperatingMode {
	case AddOnly:
		operations = []func() error{
			s.perform(ctx, "add", things, s.getThingsToAdd, adapter.Add),
		}
	case RemoveOnly:
		operations = []func() error{
			s.perform(ctx, "remove", things, s.getThingsToRemove, adapter.Remove),
		}
	case RemoveAdd:
		operations = []func() error{
			s.perform(ctx, "remove", things, s.getThingsToRemove, adapter.Remove),
			s.perform(ctx, "add", things, s.getThingsToAdd, adapter.Add),
		}
	case AddRemove:
		operations = []func() error{
			s.perform(ctx, "add", things, s.getThingsToAdd, adapter.Add),
			s.perform(ctx, "remove", things, s.getThingsToRemove, adapter.Remove),
		}
	}

	for _, fn := range operations {
		err = fn()
		if err != nil {
			return fmt.Errorf("sync.syncwith.execute -> %w", err)
		}
	}

	s.Logger.Println("Finished sync")

	return nil
}
