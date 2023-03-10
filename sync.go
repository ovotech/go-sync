package gosync

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

// Ensure Sync fully satisfies the Service interface.
var _ Service = &Sync{}

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
	DryRun bool // DryRun mode calculates membership, but doesn't add or remove.

	operatingMode OperatingMode // Change the order of Sync's operation. Default is RemoveAdd.
	caseSensitive bool          // caseSensitive sets if Go Sync is case-sensitive. Default is true.
	/*
		maximumChanges sets the maximum number of allowed changes per add/remove operation. It is not a cumulative
		total, and the number only applies to each distinct operation.

		For example:

		Setting this value to 3 means that a maximum of 3 things can be added AND removed from a destination (total 6)
		changes before Sync returns an ErrTooManyChanges error.

		Default is NoChangeLimit (or -1).
	*/
	maximumChanges int
	source         Adapter         // The source adapter.
	cache          map[string]bool // cache prevents polling the source more than once.
	logger         *log.Logger
}

// generateHashMap takes a list of strings and returns a hashed map of { item => true }.
func (s *Sync) generateHashMap(i []string) map[string]bool {
	out := map[string]bool{}

	for _, str := range i {
		if s.caseSensitive {
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
		s.logger.Println("Getting things from source adapter")

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
		s.logger.Printf("Processing things to %s\n", action)

		thingsToChange := diffFn(things)

		// If the changes exceed the maximum change limit, fail with the ErrTooManyChanges error.
		if len(thingsToChange) > s.maximumChanges && s.maximumChanges != NoChangeLimit {
			return fmt.Errorf("%s(%v) -> %w(%v)", action, thingsToChange, ErrTooManyChanges, s.maximumChanges)
		}

		if s.DryRun {
			s.logger.Printf("Would %s %s, but running in dry run mode", action, thingsToChange)

			return nil
		}

		if len(thingsToChange) == 0 {
			return nil
		}

		s.logger.Printf("%s: %s", action, thingsToChange)

		err := executeFn(ctx, thingsToChange)
		if err != nil {
			return fmt.Errorf("%s(%v) -> %w", action, things, err)
		}

		return nil
	}
}

// SyncWith synchronises the destination service with the source service, adding & removing things as necessary.
func (s *Sync) SyncWith(ctx context.Context, adapter Adapter) error {
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

	s.logger.Printf("Running in %s operating mode", s.operatingMode)

	operations := make([]func() error, 0, 2) //nolint:gomnd

	switch s.operatingMode {
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

	s.logger.Println("Finished sync")

	return nil
}

// SetOperatingMode sets a custom operating mode.
func SetOperatingMode(mode OperatingMode) func(*Sync) {
	return func(s *Sync) {
		s.operatingMode = mode
	}
}

// SetCaseSensitive can configure Go Sync's case sensitivity when comparing things.
func SetCaseSensitive(caseSensitive bool) func(*Sync) {
	return func(s *Sync) {
		s.caseSensitive = caseSensitive
	}
}

// SetMaximumChanges sets a maximum number of changes that Sync will allow before returning an ErrTooManyChanges error.
func SetMaximumChanges(maximumChanges int) func(*Sync) {
	return func(s *Sync) {
		s.maximumChanges = maximumChanges
	}
}

// New creates a new gosync.Sync service.
func New(source Adapter, configFns ...func(*Sync)) *Sync {
	sync := &Sync{
		DryRun: false,

		operatingMode:  RemoveAdd,
		caseSensitive:  true,
		maximumChanges: NoChangeLimit,
		source:         source,
		cache:          make(map[string]bool),

		logger: log.New(os.Stderr, "[go-sync/sync] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range configFns {
		fn(sync)
	}

	return sync
}
