package config

import (
	"context"
	"fmt"

	"github.com/ovotech/go-sync/internal/gosync"
)

type Options struct {
	OperatingMode  *string `yaml:"operatingMode"`
	MaximumChanges *uint16 `yaml:"maximumChanges"`
	CaseSensitive  *bool   `yaml:"caseSensitive"`
}

type Job struct {
	Options *Options   `yaml:"options,omitempty"`
	Adapter *Adapter   `validate:"required"                     yaml:"adapter"`
	With    []*Adapter `validate:"required,gte=1,dive,required" yaml:"with"`
}

func (s *Job) syncWith(ctx context.Context, sync *gosync.Sync, plugin *Adapter) error {
	destination, kill, err := plugin.ToAdapter(ctx)
	if err != nil {
		return fmt.Errorf(
			"src(%s, %v) -> dest(%s, %v).load -> %w",
			s.Adapter.Plugin,
			s.Adapter.Config,
			plugin.Plugin,
			plugin.Config,
			err,
		)
	}

	defer kill()

	err = sync.SyncWith(ctx, destination)
	if err != nil {
		return fmt.Errorf(
			"src(%s, %v) -> dest(%s, %v).syncwith -> %w",
			s.Adapter.Plugin,
			s.Adapter.Config,
			plugin.Plugin,
			plugin.Config,
			err,
		)
	}

	return nil
}

// Run this job.
func (s *Job) Run(ctx context.Context, dryRun bool) []error {
	errors := make([]error, 0)

	source, kill, err := s.Adapter.ToAdapter(ctx)
	if err != nil {
		return []error{fmt.Errorf("src(%s, %v) -> %w", s.Adapter.Plugin, s.Adapter.Config, err)}
	}

	defer kill()

	sync := gosync.New(source, func(sync *gosync.Sync) {
		sync.DryRun = dryRun
	})

	if s.Options != nil {
		if s.Options.MaximumChanges != nil {
			sync.MaximumChanges = int(*s.Options.MaximumChanges)
		}

		if s.Options.CaseSensitive != nil {
			sync.CaseSensitive = *s.Options.CaseSensitive
		}

		if s.Options.OperatingMode != nil {
			sync.OperatingMode = gosync.OperatingMode(*s.Options.OperatingMode)
		}
	}

	for _, dest := range s.With {
		err := s.syncWith(ctx, sync, dest)
		if err != nil {
			errors = append(errors, err)

			break
		}
	}

	return errors
}
