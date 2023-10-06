// Package config handles loading configuration files in that contain the specification used for Go Sync.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	gosyncerrors "github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/types"
)

type Options struct {
	OperatingMode  string `yaml:"operatingMode"`
	MaximumChanges uint16 `yaml:"maximumChanges"`
}

type Adapter struct {
	Plugin string                     `validate:"required,file" yaml:"plugin"`
	Name   string                     `validate:"required"      yaml:"name"`
	Config map[types.ConfigKey]string `yaml:"config,omitempty"`
}

type Source struct {
	Options *Options   `yaml:"options,omitempty"`
	Adapter *Adapter   `validate:"required"                     yaml:"adapter"`
	With    []*Adapter `validate:"required,gte=1,dive,required" yaml:"with"`
}

type Config struct {
	Version string   `validate:"required,eq=1"                yaml:"version"`
	Sync    []Source `validate:"required,gte=1,dive,required" yaml:"sync"`
}

// GetEnvironmentVariables intended for use with Go Sync.
func GetEnvironmentVariables() map[types.ConfigKey]string {
	vars := os.Environ()
	out := make(map[string]string)

	for _, envVar := range vars {
		if strings.HasPrefix(envVar, "GOSYNC_") {
			if key, value, ok := strings.Cut(envVar, "="); ok {
				key, _ = strings.CutPrefix(key, "GOSYNC_")

				out[key] = value
			}
		}
	}

	return out
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
