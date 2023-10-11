// Package config handles loading configuration files in that contain the specification used for Go Sync.
package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	gosyncerrors "github.com/ovotech/go-sync/pkg/errors"
)

type Config struct {
	Version string `validate:"required,eq=1"                yaml:"version"`
	Jobs    []*Job `validate:"required,gte=1,dive,required" yaml:"jobs"`
}

// Run the configuration, and execute all the jobs.
func (c *Config) Run(ctx context.Context, dryRun bool) {
	for _, job := range c.Jobs {
		errs := job.Run(ctx, dryRun)
		if len(errs) > 0 {
			log.Println("Encountered the following errors:")

			for _, err := range errs {
				log.Println(err)
			}
		}
	}
}

// Load a configuration file in, parse and validate it into a Config struct.
func Load(path string) (*Config, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Var(path, "file"); err != nil {
		return nil, fmt.Errorf("%w: %s", gosyncerrors.ErrDoesNotExist, path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open -> %w", err)
	}

	defer file.Close()

	var config Config

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, gosyncerrors.ErrInvalidConfig
	}

	if err := validate.Struct(&config); err != nil {
		var validationErrors validator.ValidationErrors

		if ok := errors.As(err, &validationErrors); ok {
			if len(validationErrors) == 0 {
				return nil, gosyncerrors.ErrInvalidConfig
			}

			return nil, fmt.Errorf("%w: %s", gosyncerrors.ErrInvalidConfig, validationErrors[0].Field())
		}

		return nil, fmt.Errorf("validate -> %w", err)
	}

	return &config, nil
}
