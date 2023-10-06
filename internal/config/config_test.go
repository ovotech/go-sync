package config_test

import (
	"github.com/ovotech/go-sync/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ovotech/go-sync/internal/config"
	gosyncerrors "github.com/ovotech/go-sync/pkg/errors"
)

//nolint:funlen
func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		expected := &config.Config{
			Version: "1",
			Sync: []config.Source{
				{
					Adapter: &config.Adapter{
						Plugin: "./fixtures/valid.yml",
						Name:   "some_adapter",
						Config: map[types.ConfigKey]string{
							"foo": "bar",
						},
					},
					With: []*config.Adapter{
						{
							Plugin: "./fixtures/valid.yml",
							Name:   "adapter_one",
							Config: map[types.ConfigKey]string{
								"some": "config",
							},
						},
						{
							Plugin: "./fixtures/valid.yml",
							Name:   "adapter_two",
							Config: map[types.ConfigKey]string{
								"more": "config",
							},
						},
					},
				},
				{
					Adapter: &config.Adapter{
						Plugin: "./fixtures/valid.yml",
						Name:   "a_different_adapter",
					},
					Options: &config.Options{
						OperatingMode: "add",
					},
					With: []*config.Adapter{
						{
							Plugin: "./fixtures/valid.yml",
							Name:   "adapter_one",
							Config: map[string]string{
								"some": "config",
							},
						},
					},
				},
			},
		}

		actual, err := config.Load("fixtures/valid.yml")

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		_, err := config.Load("fixtures/invalid.yml")

		assert.ErrorIs(t, err, gosyncerrors.ErrInvalidConfig)
	})

	t.Run("missing", func(t *testing.T) {
		t.Parallel()

		_, err := config.Load("fixtures/missing.yml")

		assert.ErrorIs(t, err, gosyncerrors.ErrInvalidConfig)
	})

	t.Run("non-existent", func(t *testing.T) {
		t.Parallel()

		_, err := config.Load("fixtures/__NOT_A_FILE__.yml")

		assert.ErrorIs(t, err, gosyncerrors.ErrDoesNotExist)
	})
}

func TestGetEnvironmentVariables(t *testing.T) { //nolint:paralleltest
	expected := map[string]string{
		"TEST_VAR": "true",
	}

	t.Setenv("GOSYNC_TEST_VAR", "true")

	actual := config.GetEnvironmentVariables()

	assert.Equal(t, expected, actual)
}
