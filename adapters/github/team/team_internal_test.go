package team

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/go-github/v47/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/github/discovery/saml"
)

func TestTeam_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := newMockIGitHubTeam(t)
	discovery := NewMockGitHubDiscovery(t)

	adapter := &Team{
		teams:     gitHubClient,
		discovery: discovery,
		org:       "org",
		slug:      "slug",
		Logger:    log.New(os.Stdout, "", log.LstdFlags),
	}

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

	adapter := &Team{
		teams:     gitHubClient,
		discovery: discovery,
		org:       "org",
		slug:      "slug",
		Logger:    log.New(os.Stdout, "", log.LstdFlags),
	}

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

	adapter := &Team{
		teams:     gitHubClient,
		discovery: discovery,
		org:       "org",
		slug:      "slug",
		cache:     map[string]string{"foo@email": "foo", "bar@email": "bar"},
		Logger:    log.New(os.Stdout, "", log.LstdFlags),
	}

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
		assert.Equal(t, "org", adapter.org)
		assert.Equal(t, "slug", adapter.slug)
		assert.IsType(t, &saml.Saml{}, adapter.discovery)
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
	})

	t.Run("invalid config", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			GitHubToken:        "token",
			GitHubOrg:          "org",
			TeamSlug:           "slug",
			DiscoveryMechanism: "foo",
		})

		assert.ErrorIs(t, err, gosync.ErrMissingConfig)
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(os.Stderr, "custom logger", log.LstdFlags)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			GitHubToken:        "token",
			GitHubOrg:          "org",
			TeamSlug:           "slug",
			DiscoveryMechanism: "saml",
		}, WithLogger(logger))

		assert.NoError(t, err)
		assert.Equal(t, logger, adapter.Logger)
	})

	t.Run("with client", func(t *testing.T) {
		t.Parallel()

		client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "token"},
		)))

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			GitHubToken:        "token",
			GitHubOrg:          "org",
			TeamSlug:           "slug",
			DiscoveryMechanism: "saml",
		}, WithClient(client))

		assert.NoError(t, err)
		assert.Equal(t, client.Teams, adapter.teams)
	})

	t.Run("with discovery service", func(t *testing.T) {
		t.Parallel()

		mockDiscovery := NewMockGitHubDiscovery(t)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			GitHubToken: "token",
			GitHubOrg:   "org",
			TeamSlug:    "slug",
		}, WithDiscoveryService(mockDiscovery))

		assert.NoError(t, err)
		assert.Equal(t, mockDiscovery, adapter.discovery)
	})
}
