package github

import (
	"testing"

	"github.com/google/go-github/v47/github"
	"github.com/ovotech/go_sync/test/httpclient"
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

func TestNewGitHubTeamService(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	gitHubClient := github.NewClient(httpClient)
	discovery := new(mockDiscovery)
	repo := NewTeamSyncService(gitHubClient, discovery, "org", "slug")

	assert.IsType(t, new(Team), repo)
	assert.Zero(t, mockRoundTrip.Calls)
}

func TestGitHubTeam_Remove_Error(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGitHubRequest(
		"DELETE",
		"/orgs/org/teams/slug/memberships/foo",
		403,
		"",
	)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubClient := github.NewClient(httpClient)
	repo := NewTeamSyncService(gitHubClient, nil, "org", "slug")
	repo.cache = map[string]string{"foo@test.test": "foo"}

	success, failure, err := repo.Remove("foo@test.test")

	assert.Empty(t, success)
	assert.NotEmpty(t, failure)
	assert.NoError(t, err)
}
