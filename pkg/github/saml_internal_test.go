package github

import (
	"testing"

	"github.com/ovotech/go-sync/test/httpclient"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

func TestNewGitHubSaml(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	gitHubV4Client := githubv4.NewClient(httpClient)
	discovery := NewSamlDiscoveryService(gitHubV4Client, "org")

	assert.IsType(t, new(Saml), discovery)
	assert.Equal(t, "org", discovery.org)
	assert.Equal(t, gitHubV4Client, discovery.client)
	assert.Zero(t, mockRoundTrip.Calls)
}
