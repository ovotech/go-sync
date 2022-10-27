package team

import (
	"context"
	"testing"

	"github.com/google/go-github/v47/github"
	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/github/discovery/saml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	discovery := NewMockGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")

	assert.Equal(t, "org", adapter.org)
	assert.Equal(t, "slug", adapter.slug)
}

func TestTeam_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := newMockIGitHubTeam(t)
	discovery := NewMockGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")
	adapter.teams = gitHubClient

	gitHubClient.
		EXPECT().
		ListTeamMembersBySlug(ctx, "org", "slug", &github.TeamListTeamMembersOptions{}).
		Return([]*github.User{
			{Login: github.String("foo")},
		}, &github.Response{NextPage: 1}, nil)
	gitHubClient.
		EXPECT().
		ListTeamMembersBySlug(ctx, "org", "slug", &github.TeamListTeamMembersOptions{
			ListOptions: github.ListOptions{Page: 1},
		}).
		Return([]*github.User{
			{Login: github.String("bar")},
		}, &github.Response{NextPage: 0}, nil)
	discovery.EXPECT().GetEmailFromUsername(ctx, []string{"foo"}).Return([]string{"foo@email"}, nil)
	discovery.EXPECT().GetEmailFromUsername(ctx, []string{"bar"}).Return([]string{"bar@email"}, nil)

	users, err := adapter.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, users, []string{"foo@email", "bar@email"})
	assert.Equal(t, map[string]string{"foo@email": "foo", "bar@email": "bar"}, adapter.cache)
}

func TestTeam_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := newMockIGitHubTeam(t)
	discovery := NewMockGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")
	adapter.teams = gitHubClient

	discovery.EXPECT().GetUsernameFromEmail(ctx, []string{"fizz@email", "buzz@email"}).
		Maybe().Return([]string{"fizz", "buzz"}, nil)
	discovery.EXPECT().GetUsernameFromEmail(ctx, []string{"fizz@email", "buzz@email"}).
		Maybe().Return([]string{"buzz", "fizz"}, nil)
	gitHubClient.EXPECT().AddTeamMembershipBySlug(ctx, "org", "slug", "fizz", mock.Anything).Return(nil, nil, nil)
	gitHubClient.EXPECT().AddTeamMembershipBySlug(ctx, "org", "slug", "buzz", mock.Anything).Return(nil, nil, nil)

	err := adapter.Add(ctx, []string{"fizz@email", "buzz@email"})

	assert.NoError(t, err)
}

func TestTeam_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := newMockIGitHubTeam(t)
	discovery := NewMockGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")
	adapter.teams = gitHubClient
	adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

	gitHubClient.EXPECT().RemoveTeamMembershipBySlug(ctx, "org", "slug", "foo").Return(nil, nil)
	gitHubClient.EXPECT().RemoveTeamMembershipBySlug(ctx, "org", "slug", "bar").Return(nil, nil)

	err := adapter.Remove(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
}

//nolint:funlen
func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			GitHubToken:        "token",
			GitHubOrg:          "org",
			TeamSlug:           "slug",
			DiscoveryMechanism: "saml",
		})

		assert.NoError(t, err)
		assert.IsType(t, &Team{}, adapter)
		assert.Equal(t, "org", adapter.(*Team).org)
		assert.Equal(t, "slug", adapter.(*Team).slug)
		assert.IsType(t, &saml.Saml{}, adapter.(*Team).discovery)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing token", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				GitHubOrg:          "org",
				TeamSlug:           "slug",
				DiscoveryMechanism: "saml",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, GitHubToken)
		})

		t.Run("missing org", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				GitHubToken:        "token",
				TeamSlug:           "slug",
				DiscoveryMechanism: "saml",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, GitHubOrg)
		})

		t.Run("missing slug", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				GitHubToken:        "token",
				GitHubOrg:          "org",
				DiscoveryMechanism: "saml",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, TeamSlug)
		})

		t.Run("missing discovery", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				GitHubToken: "token",
				GitHubOrg:   "org",
				TeamSlug:    "slug",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, DiscoveryMechanism)
		})
	})

	t.Run("invalid config", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			GitHubToken:        "token",
			GitHubOrg:          "org",
			TeamSlug:           "slug",
			DiscoveryMechanism: "foo",
		})

		assert.ErrorIs(t, err, gosync.ErrInvalidConfig)
		assert.ErrorContains(t, err, "foo")
	})
}
