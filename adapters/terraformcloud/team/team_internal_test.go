package team

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gosync "github.com/ovotech/go-sync"
)

func TestTeam_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	iTeamsClient := newMockITeams(t)

	adapter := &Team{
		organisation: "test",
		teams:        iTeamsClient,
		Logger:       log.New(os.Stdout, "", log.LstdFlags),
	}

	iTeamsClient.EXPECT().List(ctx, "test", &tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1},
	}).Return(&tfe.TeamList{
		Pagination: &tfe.Pagination{
			CurrentPage: 1,
			NextPage:    2,
			TotalPages:  2,
		},
		Items: []*tfe.Team{{Name: "foo"}, {Name: "bar"}},
	}, nil)

	iTeamsClient.EXPECT().List(ctx, "test", &tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 2},
	}).Return(&tfe.TeamList{
		Pagination: &tfe.Pagination{
			CurrentPage: 2,
			NextPage:    2,
			TotalPages:  2,
		},
		Items: []*tfe.Team{{Name: "fizz"}, {Name: "buzz"}},
	}, nil)

	things, err := adapter.Get(ctx)

	require.NoError(t, err)
	assert.ElementsMatch(t, things, []string{"foo", "bar", "fizz", "buzz"})
}

func TestTeam_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	iTeamsClient := newMockITeams(t)

	adapter := &Team{
		organisation: "test",
		teams:        iTeamsClient,
		Logger:       log.New(os.Stdout, "", log.LstdFlags),
	}

	foo := "foo"

	iTeamsClient.EXPECT().Create(ctx, "test", tfe.TeamCreateOptions{Name: &foo}).Return(&tfe.Team{}, nil)

	err := adapter.Add(ctx, []string{"foo"})

	require.NoError(t, err)
}

func TestTeam_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	iTeamsClient := newMockITeams(t)

	adapter := &Team{
		organisation: "test",
		teams:        iTeamsClient,
		cache:        map[string]string{"foo": "foo-id"},
		Logger:       log.New(os.Stdout, "", log.LstdFlags),
	}

	iTeamsClient.EXPECT().Delete(ctx, "foo-id").Return(nil)

	err := adapter.Remove(ctx, []string{"foo"})

	require.NoError(t, err)
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{Token: "token", Organisation: "org"})

		require.NoError(t, err)
		assert.IsType(t, &Team{}, adapter)
	})

	t.Run("missing token", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{Organisation: "org"})

		require.ErrorIs(t, err, gosync.ErrMissingConfig)
		require.ErrorContains(t, err, Token)
	})

	t.Run("missing organisation", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{Token: "token"})

		require.ErrorIs(t, err, gosync.ErrMissingConfig)
		require.ErrorContains(t, err, Organisation)
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(os.Stderr, "custom logger", log.LstdFlags)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Token:        "token",
			Organisation: "org",
		}, WithLogger(logger))

		require.NoError(t, err)
		assert.Equal(t, logger, adapter.Logger)
	})

	t.Run("with client", func(t *testing.T) {
		t.Parallel()

		client, err := tfe.NewClient(&tfe.Config{Token: "test"})
		require.NoError(t, err)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Organisation: "org",
		}, WithClient(client))

		require.NoError(t, err)
		assert.Equal(t, client.Teams, adapter.teams)
	})
}
