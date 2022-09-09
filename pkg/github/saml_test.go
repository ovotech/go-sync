package github_test

import (
	"testing"

	"github.com/ovotech/go_sync/pkg/github"
	"github.com/ovotech/go_sync/test/httpclient"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

func TestSaml_GetUsernameFromEmail(t *testing.T) { //nolint:funlen
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGraphQlRequest(
		"query($email:String!$org:String!){organization(login: $org){samlIdentityProvider{externalIdentities(userName: $email, first: 1){edges{node{user{login}}}}}}}", //nolint:lll
		map[string]string{
			"org":   "org",
			"email": "foo@test.test",
		},
		`
{
	"data": {
		"organization": {
			"samlIdentityProvider": {
				"externalIdentities": {
					"edges": [
						{
							"node": {
								"user": {
									"login": "foo"
								}
							}
						}
					]
				}
			}
		}
	}
}
`)
	mockRoundTrip.AddGraphQlRequest(
		"query($email:String!$org:String!){organization(login: $org){samlIdentityProvider{externalIdentities(userName: $email, first: 1){edges{node{user{login}}}}}}}", //nolint:lll
		map[string]string{
			"org":   "org",
			"email": "bar@test.test",
		},
		`
{
	"data": {
		"organization": {
			"samlIdentityProvider": {
				"externalIdentities": {
					"edges": [
						{
							"node": {
								"user": {
									"login": "bar"
								}
							}
						}
					]
				}
			}
		}
	}
}
`)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubV4Client := githubv4.NewClient(httpClient)
	discovery := github.NewSamlDiscoveryService(gitHubV4Client, "org")

	ids, err := discovery.GetUsernameFromEmail("foo@test.test", "bar@test.test")

	assert.ElementsMatch(t, []string{"foo", "bar"}, ids)
	assert.NoError(t, err)
}

func TestSaml_GetUsernameFromEmail_UserNotFound(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGraphQlRequest(
		"query($email:String!$org:String!){organization(login: $org){samlIdentityProvider{externalIdentities(userName: $email, first: 1){edges{node{user{login}}}}}}}", //nolint:lll
		map[string]string{
			"org":   "org",
			"email": "foo@test.test",
		},
		`
{
	"data": {
		"organization": {
			"samlIdentityProvider": {
				"externalIdentities": {
					"edges": []
				}
			}
		}
	}
}
`)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubV4Client := githubv4.NewClient(httpClient)
	discovery := github.NewSamlDiscoveryService(gitHubV4Client, "org")

	_, err := discovery.GetUsernameFromEmail("foo@test.test")

	assert.Error(t, err)
	assert.ErrorIs(t, err, github.ErrUserNotFound)
}

func TestSaml_GetEmailFromUsername(t *testing.T) { //nolint:funlen
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGraphQlRequest(
		"query($login:String!$org:String!){organization(login: $org){samlIdentityProvider{externalIdentities(login: $login, first: 1){edges{node{samlIdentity{nameId}}}}}}}", //nolint:lll
		map[string]string{
			"org":   "org",
			"login": "foo",
		},
		`
{
	"data": {
		"organization": {
			"samlIdentityProvider": {
				"externalIdentities": {
					"edges": [
						{
							"node": {
								"samlIdentity": {
									"nameId": "foo@test.test"
								}
							}
						}
					]
				}
			}
		}
	}
}
`)
	mockRoundTrip.AddGraphQlRequest(
		"query($login:String!$org:String!){organization(login: $org){samlIdentityProvider{externalIdentities(login: $login, first: 1){edges{node{samlIdentity{nameId}}}}}}}", //nolint:lll
		map[string]string{
			"org":   "org",
			"login": "bar",
		},
		`
{
	"data": {
		"organization": {
			"samlIdentityProvider": {
				"externalIdentities": {
					"edges": [
						{
							"node": {
								"samlIdentity": {
									"nameId": "bar@test.test"
								}
							}
						}
					]
				}
			}
		}
	}
}
`)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubV4Client := githubv4.NewClient(httpClient)
	discovery := github.NewSamlDiscoveryService(gitHubV4Client, "org")

	ids, err := discovery.GetEmailFromUsername("foo", "bar")

	assert.ElementsMatch(t, []string{"foo@test.test", "bar@test.test"}, ids)
	assert.NoError(t, err)
}

func TestSaml_GetEmailFromUsername_UserNotFound(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	mockRoundTrip.AddGraphQlRequest(
		"query($login:String!$org:String!){organization(login: $org){samlIdentityProvider{externalIdentities(login: $login, first: 1){edges{node{samlIdentity{nameId}}}}}}}", //nolint:lll
		map[string]string{
			"org":   "org",
			"login": "foo",
		},
		`
{
	"data": {
		"organization": {
			"samlIdentityProvider": {
				"externalIdentities": {
					"edges": []
				}
			}
		}
	}
}
`)

	httpClient := httpclient.New(mockRoundTrip)
	gitHubV4Client := githubv4.NewClient(httpClient)
	discovery := github.NewSamlDiscoveryService(gitHubV4Client, "org")

	_, err := discovery.GetEmailFromUsername("foo")

	assert.Error(t, err)
	assert.ErrorIs(t, err, github.ErrUserNotFound)
}
