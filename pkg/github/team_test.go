package github_test

import (
	"errors"
	"testing"

	"github.com/google/go-github/v47/github"
	githubsync "github.com/ovotech/go-sync/pkg/github"
	"github.com/ovotech/go-sync/test/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDiscovery struct {
	mock.Mock
}

func (d *mockDiscovery) GetUsernameFromEmail(in ...string) ([]string, error) {
	args := d.Called(in)

	return args.Get(0).([]string), args.Error(1) //nolint:wrapcheck
}

func (d *mockDiscovery) GetEmailFromUsername(in ...string) ([]string, error) {
	args := d.Called(in)

	return args.Get(0).([]string), args.Error(1) //nolint:wrapcheck
}

func TestGitHubTeam_Get(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGitHubRequest(
		"GET",
		"/orgs/org/teams/slug/members",
		200,
		`
[
	{
		"login": "foo"
	}
]
`)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubClient := github.NewClient(httpClient)
	discovery := new(mockDiscovery)
	discovery.On("GetEmailFromUsername", []string{"foo"}).Return([]string{"foo@test.test"}, nil)
	repo := githubsync.NewTeamSyncService(gitHubClient, discovery, "org", "slug")

	emails, err := repo.Get()

	assert.Equal(t, []string{"foo@test.test"}, emails)
	assert.NoError(t, err)
}

func TestGitHubTeam_Get_Error(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGitHubRequest(
		"GET",
		"/orgs/org/teams/slug/members",
		404,
		`
{
	"message": "Not Found",
	"documentation_url": "https://docs.github.com/rest/reference/teams#list-team-members"
}
`)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubClient := github.NewClient(httpClient)
	repo := githubsync.NewTeamSyncService(gitHubClient, nil, "org", "slug")

	_, err := repo.Get()

	assert.Error(t, err)
}

func TestGitHubTeam_Add(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGitHubRequest(
		"PUT",
		"/orgs/org/teams/slug/memberships/foo",
		200,
		`
{
	"url": "foo",
	"role": "member",
	"state": "pending"
}
`)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubClient := github.NewClient(httpClient)
	discovery := new(mockDiscovery)
	discovery.On("GetUsernameFromEmail", []string{"foo@test.test"}).Return([]string{"foo"}, nil)
	repo := githubsync.NewTeamSyncService(gitHubClient, discovery, "org", "slug")

	success, failure, err := repo.Add("foo@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestGitHubTeam_Add_NoUser(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubClient := github.NewClient(httpClient)
	discovery := new(mockDiscovery)
	discovery.On("GetUsernameFromEmail", []string{"foo@test.test"}).Return([]string(nil), errors.New("user_not_found")) //nolint:lll,goerr113
	repo := githubsync.NewTeamSyncService(gitHubClient, discovery, "org", "slug")

	_, _, err := repo.Add("foo@test.test")

	assert.ErrorContains(t, err, "user_not_found")
}

func TestGitHubTeam_Remove_NoCache(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGitHubRequest(
		"GET",
		"/orgs/org/teams/slug/members",
		200,
		`
[
	{
		"login": "foo"
	}
]
`)
	mockRoundTrip.AddGitHubRequest(
		"DELETE",
		"/orgs/org/teams/slug/memberships/foo",
		204,
		"",
	)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubClient := github.NewClient(httpClient)
	discovery := new(mockDiscovery)
	discovery.On("GetEmailFromUsername", []string{"foo"}).Return([]string{"foo@test.test"}, nil)
	repo := githubsync.NewTeamSyncService(gitHubClient, discovery, "org", "slug")

	success, failure, err := repo.Remove("foo@test.test")

	assert.Equal(t, []string{"foo@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}
