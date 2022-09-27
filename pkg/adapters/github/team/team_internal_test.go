package team

import (
	"context"
	"testing"

	"github.com/google/go-github/v47/github"
	"github.com/ovotech/go-sync/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	discovery := mocks.NewGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")

	assert.Equal(t, "org", adapter.org)
	assert.Equal(t, "slug", adapter.slug)
}

func TestTeam_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := mocks.NewIGitHubTeam(t)
	discovery := mocks.NewGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")
	adapter.client = gitHubClient

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

	gitHubClient := mocks.NewIGitHubTeam(t)
	discovery := mocks.NewGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")
	adapter.client = gitHubClient

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

	gitHubClient := mocks.NewIGitHubTeam(t)
	discovery := mocks.NewGitHubDiscovery(t)
	adapter := New(&github.Client{}, discovery, "org", "slug")
	adapter.client = gitHubClient
	adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

	gitHubClient.EXPECT().RemoveTeamMembershipBySlug(ctx, "org", "slug", "foo").Return(nil, nil)
	gitHubClient.EXPECT().RemoveTeamMembershipBySlug(ctx, "org", "slug", "bar").Return(nil, nil)

	err := adapter.Remove(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
}
