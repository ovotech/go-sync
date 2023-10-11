package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gosyncerrors "github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/types"
)

//nolint:funlen
func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		add := "add"

		expected := &Config{
			Version: "1",
			Jobs: []*Job{
				{
					Adapter: &Adapter{
						Plugin: "./fixtures/valid.yml",
						Config: map[types.ConfigKey]string{
							"foo": "bar",
						},
					},
					With: []*Adapter{
						{
							Plugin: "./fixtures/valid.yml",
							Config: map[types.ConfigKey]string{
								"some": "config",
							},
						},
						{
							Plugin: "./fixtures/valid.yml",
							Config: map[types.ConfigKey]string{
								"more": "config",
							},
						},
					},
				},
				{
					Adapter: &Adapter{
						Plugin: "./fixtures/valid.yml",
					},
					Options: &Options{
						OperatingMode: &add,
					},
					With: []*Adapter{
						{
							Plugin: "./fixtures/valid.yml",
							Config: map[string]string{
								"some": "config",
							},
						},
					},
				},
			},
		}

		actual, err := Load("fixtures/valid.yml")

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		_, err := Load("fixtures/invalid.yml")

		assert.ErrorIs(t, err, gosyncerrors.ErrInvalidConfig)
	})

	t.Run("missing", func(t *testing.T) {
		t.Parallel()

		_, err := Load("fixtures/missing.yml")

		assert.ErrorIs(t, err, gosyncerrors.ErrInvalidConfig)
	})

	t.Run("non-existent", func(t *testing.T) {
		t.Parallel()

		_, err := Load("fixtures/__NOT_A_FILE__.yml")

		assert.ErrorIs(t, err, gosyncerrors.ErrDoesNotExist)
	})
}

func TestGetEnvironmentVariables(t *testing.T) { //nolint:paralleltest
	expected := map[string]string{
		"TEST_VAR": "true",
	}

	t.Setenv("GOSYNC_TEST_VAR", "true")

	actual := getEnvironmentVariables()

	assert.Equal(t, expected, actual)
}
