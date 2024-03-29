package user

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

func TestUser_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockTeams := newMockITeams(t)

	adapter := &User{
		organisation: "org",
		team:         "team",
		teams:        mockTeams,
		Logger:       log.New(os.Stdout, "", log.LstdFlags),
	}

	mockTeams.EXPECT().List(ctx, "org", &tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 0,
		},
		Include: []tfe.TeamIncludeOpt{tfe.TeamOrganizationMemberships},
		Names:   []string{"team"},
	}).Return(&tfe.TeamList{
		Items: []*tfe.Team{
			{
				OrganizationMemberships: []*tfe.OrganizationMembership{
					{Email: "foo@email"},
					{Email: "bar@email"},
				},
			},
		},
	}, nil)

	things, err := adapter.Get(ctx)

	require.NoError(t, err)
	assert.ElementsMatch(t, things, []string{"foo@email", "bar@email"})
}

func TestUser_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockTeams := newMockITeams(t)
	mockOrgMembership := newMockIOrganizationMemberships(t)
	mockTeamMembers := newMockITeamMembers(t)

	adapter := &User{
		organisation:            "org",
		team:                    "team",
		teams:                   mockTeams,
		organizationMemberships: mockOrgMembership,
		teamMembers:             mockTeamMembers,
		Logger:                  log.New(os.Stdout, "", log.LstdFlags),
	}

	// Mock a first page of responses from the API.
	mockOrgMembership.EXPECT().List(ctx, "org", &tfe.OrganizationMembershipListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1},
		Emails:      []string{"foo@email", "bar@email"},
	}).Return(&tfe.OrganizationMembershipList{
		Pagination: &tfe.Pagination{
			CurrentPage: 1,
			NextPage:    2,
			TotalPages:  2,
		},
		Items: []*tfe.OrganizationMembership{
			{Email: "foo@email", ID: "foo"},
		},
	}, nil)

	// Mock a second page of responses from the API.
	mockOrgMembership.EXPECT().List(ctx, "org", &tfe.OrganizationMembershipListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 2},
		Emails:      []string{"foo@email", "bar@email"},
	}).Return(&tfe.OrganizationMembershipList{
		Pagination: &tfe.Pagination{
			CurrentPage: 2,
			NextPage:    2,
			TotalPages:  2,
		},
		Items: []*tfe.OrganizationMembership{
			{Email: "bar@email", ID: "bar"},
		},
	}, nil)

	// Mock converting the friendly team name into an ID.
	mockTeams.EXPECT().List(ctx, "org", &tfe.TeamListOptions{
		Names: []string{"team"},
	}).Return(&tfe.TeamList{
		Items: []*tfe.Team{{ID: "team-id"}},
	}, nil)

	mockTeamMembers.EXPECT().Add(ctx, "team-id", tfe.TeamMemberAddOptions{
		OrganizationMembershipIDs: []string{"foo", "bar"},
	}).Return(nil)

	err := adapter.Add(ctx, []string{"foo@email", "bar@email"})

	require.NoError(t, err)
}

func TestUser_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockTeams := newMockITeams(t)
	mockOrgMembership := newMockIOrganizationMemberships(t)
	mockTeamMembers := newMockITeamMembers(t)

	adapter := &User{
		organisation:            "org",
		team:                    "team",
		teams:                   mockTeams,
		organizationMemberships: mockOrgMembership,
		teamMembers:             mockTeamMembers,
		Logger:                  log.New(os.Stdout, "", log.LstdFlags),
	}

	mockOrgMembership.EXPECT().List(ctx, "org", &tfe.OrganizationMembershipListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1},
		Emails:      []string{"foo@email"},
	}).Return(&tfe.OrganizationMembershipList{
		Pagination: &tfe.Pagination{
			CurrentPage: 1,
			NextPage:    1,
			TotalPages:  1,
		},
		Items: []*tfe.OrganizationMembership{
			{Email: "foo@email", ID: "foo"},
		},
	}, nil)

	// Mock converting the friendly team name into an ID.
	mockTeams.EXPECT().List(ctx, "org", &tfe.TeamListOptions{
		Names: []string{"team"},
	}).Return(&tfe.TeamList{
		Items: []*tfe.Team{{ID: "team-id"}},
	}, nil)

	mockTeamMembers.EXPECT().Remove(ctx, "team-id", tfe.TeamMemberRemoveOptions{
		OrganizationMembershipIDs: []string{"foo"},
	}).Return(nil)

	err := adapter.Remove(ctx, []string{"foo@email"})

	require.NoError(t, err)
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Token:        "token",
			Organisation: "org",
			Team:         "team",
		})

		require.NoError(t, err)
		assert.IsType(t, &User{}, adapter)
	})

	t.Run("missing token", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			Organisation: "org",
			Team:         "team",
		})

		require.ErrorIs(t, err, gosync.ErrMissingConfig)
		require.ErrorContains(t, err, Token)
	})

	t.Run("missing organisation", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			Token: "token",
			Team:  "team",
		})

		require.ErrorIs(t, err, gosync.ErrMissingConfig)
		require.ErrorContains(t, err, Organisation)
	})

	t.Run("missing team", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			Token:        "token",
			Organisation: "org",
		})

		require.ErrorIs(t, err, gosync.ErrMissingConfig)
		require.ErrorContains(t, err, Team)
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(os.Stderr, "custom logger", log.LstdFlags)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Token:        "token",
			Organisation: "org",
			Team:         "team",
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
			Team:         "team",
		}, WithClient(client))

		require.NoError(t, err)
		assert.Equal(t, client.Teams, adapter.teams)
		assert.Equal(t, client.TeamMembers, adapter.teamMembers)
		assert.Equal(t, client.OrganizationMemberships, adapter.organizationMemberships)
	})
}
