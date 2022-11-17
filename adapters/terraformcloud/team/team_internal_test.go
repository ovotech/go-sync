package team

import (
	"context"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	gosync "github.com/ovotech/go-sync"
)

func TestNew(t *testing.T) {
	t.Parallel()
}

func TestTeam_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	iTeamsClient := newMockITeams(t)
	adapter := New(&tfe.Client{}, "test")
	adapter.teams = iTeamsClient

	iTeamsClient.EXPECT().List(ctx, "test", mock.Anything).Return(&tfe.TeamList{
		Items: []*tfe.Team{{Name: "foo"}, {Name: "bar"}},
	}, nil)

	things, err := adapter.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, things, []string{"foo", "bar"})
}

func TestTeam_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	iTeamsClient := newMockITeams(t)
	adapter := New(&tfe.Client{}, "test")
	adapter.teams = iTeamsClient

	foo := "foo"

	iTeamsClient.EXPECT().Create(ctx, "test", tfe.TeamCreateOptions{Name: &foo}).Return(&tfe.Team{}, nil)

	err := adapter.Add(ctx, []string{"foo"})

	assert.NoError(t, err)
}

func TestTeam_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	iTeamsClient := newMockITeams(t)
	adapter := New(&tfe.Client{}, "test")
	adapter.teams = iTeamsClient
	adapter.cache = map[string]string{"foo": "foo-id"}

	iTeamsClient.EXPECT().Delete(ctx, "foo-id").Return(nil)

	err := adapter.Remove(ctx, []string{"foo"})

	assert.NoError(t, err)
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{Token: "token", Organisation: "org"})

		assert.NoError(t, err)
		assert.IsType(t, &Team{}, adapter)
	})

	t.Run("missing token", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{Organisation: "org"})

		assert.ErrorIs(t, err, gosync.ErrMissingConfig)
		assert.ErrorContains(t, err, Token)
	})

	t.Run("missing organisation", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{Token: "token"})

		assert.ErrorIs(t, err, gosync.ErrMissingConfig)
		assert.ErrorContains(t, err, Organisation)
	})
}
