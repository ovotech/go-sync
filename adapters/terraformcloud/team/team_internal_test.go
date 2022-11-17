package team

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	gosync "github.com/ovotech/go-sync"
)

func TestNew(t *testing.T) {
	t.Parallel()
}

func TestTeam_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	adapter := New()
	things, err := adapter.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, things, []string{})
}

func TestTeam_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	adapter := New()
	err := adapter.Add(ctx, []string{"foo"})

	assert.NoError(t, err)
}

func TestTeam_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	adapter := New()
	err := adapter.Remove(ctx, []string{"bar"})

	assert.NoError(t, err)
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{})

		assert.NoError(t, err)
		assert.IsType(t, &Team{}, adapter)
	})
}
