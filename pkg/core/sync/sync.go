package sync

import (
	"fmt"

	"github.com/ovotech/go-sync/pkg/core/ports"
)

// action-able function, such as adding or removing users.
type action func(...string) (success []string, failure []error, err error)

type Sync struct {
	source ports.Service   // The source of truth whose membership will be synced with other services.
	cache  map[string]bool // A cache of usernames to prevent polling the source too often.
	dryRun bool            // Flag to indicate if running in dryRun mode.
	add    bool            // Flag to indicate if Sync should add accounts.
	remove bool            // Flag to indicate if Sync should remove accounts.
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
func (s *Sync) generateCache() error {
	if len(s.cache) == 0 {
		accounts, err := s.source.Get()
		if err != nil {
			return fmt.Errorf("get -> %w", err)
		}

		s.cache = generateHashMap(accounts)
	}

	return nil
}

// New creates a new Sync service.
func New(source ports.Service) *Sync {
	return &Sync{
		source: source,
		cache:  map[string]bool{},
		dryRun: false,
		add:    true,
		remove: true,
	}
}

// SetDryRun can be used to enable dry-run mode (default: disabled).
// In dry-run mode, accounts will be enumerated but will not add/remove accounts.
func (s *Sync) SetDryRun(dryRun bool) {
	s.dryRun = dryRun
}

// SetAddRemove can be used to selectively enable/disable the add or remove functionality.
func (s *Sync) SetAddRemove(add bool, remove bool) {
	s.add = add
	s.remove = remove
}

// perform runs the diff func, and then actions it to return the response.
func (s *Sync) perform(accounts []string, diff func([]string) []string, action action) ([]string, []error, error) {
	accountsToAction := diff(accounts)
	if s.dryRun {
		// If running in dry-run mode, return the diff as success, but don't action the change.
		return accountsToAction, nil, nil
	}

	// If nothing needs to happen, then just return here with empty values.
	if len(accountsToAction) == 0 {
		return []string{}, []error{}, nil
	}

	return action(accountsToAction...)
}

// SyncWith synchronises the requested service with the source service, adding & removing members.
// Returns a list of successful Sync, a list of errors, and a general error for critical failures.
func (s *Sync) SyncWith(service ports.Service) ([]string, []error, error) {
	var (
		success []string
		failure []error
	)

	// Call to populate the cache from the source service.
	if err := s.generateCache(); err != nil {
		return nil, nil, fmt.Errorf("sync.SyncWith.generateCache -> %w", err)
	}

	accounts, err := service.Get()
	if err != nil {
		return nil, nil, fmt.Errorf("sync.SyncWith.Get -> %w", err)
	}

	if s.remove {
		performSuccess, performFailure, err := s.perform(accounts, s.getAccountsToRemove, service.Remove)
		if err != nil {
			return nil, nil, fmt.Errorf("sync.SyncWith.Remove -> %w", err)
		}

		success = append(success, performSuccess...)
		failure = append(failure, performFailure...)
	}

	if s.add {
		performSuccess, performFailure, err := s.perform(accounts, s.getAccountsToAdd, service.Add)
		if err != nil {
			return nil, nil, fmt.Errorf("sync.SyncWith.Add -> %w", err)
		}

		success = append(success, performSuccess...)
		failure = append(failure, performFailure...)
	}

	return success, failure, nil
}
