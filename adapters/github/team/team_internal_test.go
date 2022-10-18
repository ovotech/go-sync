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
	adapter := New(&github.Client{}, discovery, "Org", "Slug")

	assert.Equal(t, "Org", adapter.org)
	assert.Equal(t, "Slug", adapter.slug)
}

func TestTeam_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := newMockIGitHubTeam(t)
	discovery := NewMockGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "Org", "Slug")
	adapter.teams = gitHubClient

	gitHubClient.
		EXPECT().
		ListTeamMembersBySlug(ctx, "Org", "Slug", &github.TeamListTeamMembersOptions{}).
		Return([]*github.User{
			{Login: github.String("foo")},
		}, &github.Response{NextPage: 1}, nil)
	gitHubClient.
		EXPECT().
		ListTeamMembersBySlug(ctx, "Org", "Slug", &github.TeamListTeamMembersOptions{
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
	adapter := New(&github.Client{}, discovery, "Org", "Slug")
	adapter.teams = gitHubClient

	discovery.EXPECT().GetUsernameFromEmail(ctx, []string{"fizz@email", "buzz@email"}).
		Maybe().Return([]string{"fizz", "buzz"}, nil)
	discovery.EXPECT().GetUsernameFromEmail(ctx, []string{"fizz@email", "buzz@email"}).
		Maybe().Return([]string{"buzz", "fizz"}, nil)
	gitHubClient.EXPECT().AddTeamMembershipBySlug(ctx, "Org", "Slug", "fizz", mock.Anything).Return(nil, nil, nil)
	gitHubClient.EXPECT().AddTeamMembershipBySlug(ctx, "Org", "Slug", "buzz", mock.Anything).Return(nil, nil, nil)

	err := adapter.Add(ctx, []string{"fizz@email", "buzz@email"})

	assert.NoError(t, err)
}

func TestTeam_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := newMockIGitHubTeam(t)
	discovery := NewMockGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "Org", "Slug")
	adapter.teams = gitHubClient
	adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

	gitHubClient.EXPECT().RemoveTeamMembershipBySlug(ctx, "Org", "Slug", "foo").Return(nil, nil)
	gitHubClient.EXPECT().RemoveTeamMembershipBySlug(ctx, "Org", "Slug", "bar").Return(nil, nil)

	err := adapter.Remove(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
}

//nolint:funlen
func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(map[InitKey]string{
			GitHubToken:     "token",
			GitHubOrg:       "org",
			GitHubTeamSlug:  "slug",
			GitHubDiscovery: "saml",
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

			_, err := Init(map[InitKey]string{
				GitHubOrg:       "org",
				GitHubTeamSlug:  "slug",
				GitHubDiscovery: "saml",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, GitHubToken)
		})

		t.Run("missing org", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[InitKey]string{
				GitHubToken:     "token",
				GitHubTeamSlug:  "slug",
				GitHubDiscovery: "saml",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, GitHubOrg)
		})

		t.Run("missing slug", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[InitKey]string{
				GitHubToken:     "token",
				GitHubOrg:       "org",
				GitHubDiscovery: "saml",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, GitHubTeamSlug)
		})

		t.Run("missing discovery", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[InitKey]string{
				GitHubToken:    "token",
				GitHubOrg:      "org",
				GitHubTeamSlug: "slug",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, GitHubDiscovery)
		})
	})

	t.Run("invalid config", func(t *testing.T) {
		t.Parallel()

		_, err := Init(map[InitKey]string{
			GitHubToken:     "token",
			GitHubOrg:       "org",
			GitHubTeamSlug:  "slug",
			GitHubDiscovery: "foo",
		})

		assert.ErrorIs(t, err, gosync.ErrInvalidConfig)
		assert.ErrorContains(t, err, "foo")
	})
}
