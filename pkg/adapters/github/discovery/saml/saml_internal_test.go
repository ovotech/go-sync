package saml

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ovotech/go-sync/internal/mocks"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	gitHubClient := mocks.NewIGitHubV4Saml(t)
	discovery := New(nil, "test")
	discovery.client = gitHubClient

	assert.Equal(t, "test", discovery.org)
	assert.Zero(t, gitHubClient.Calls)
}

func TestSaml_GetEmailFromUsername(t *testing.T) { //nolint:dupl
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := mocks.NewIGitHubV4Saml(t)
	discovery := New(nil, "test")
	discovery.client = gitHubClient

	queryFn := func(name string) func(ctx context.Context, q interface{}, variables map[string]interface{}) {
		return func(ctx context.Context, query interface{}, variables map[string]interface{}) {
			resp := fmt.Sprintf(`
{
  "Organization": {
    "SamlIdentityProvider": {
      "ExternalIdentities": {
        "Edges": [
          {
            "Node": {
              "SamlIdentity": {
                "NameID": "%s"
              }
            }
          }
        ]
      }
    }
  }
}
`, name)
			arg := query.(*emailQuery)
			err := json.Unmarshal([]byte(resp), arg)

			assert.NoError(t, err)
			assert.Equal(t, name,
				arg.Organization.SamlIdentityProvider.ExternalIdentities.Edges[0].Node.SamlIdentity.NameID,
			)
		}
	}

	gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
		"login": githubv4.String("foo"), "org": githubv4.String("test")},
	).Run(queryFn("foo@email")).
		Return(nil)
	gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
		"login": githubv4.String("bar"), "org": githubv4.String("test")},
	).Run(queryFn("bar@email")).
		Return(nil)

	usernames, err := discovery.GetEmailFromUsername(ctx, []string{"foo", "bar"})

	assert.NoError(t, err)
	assert.ElementsMatch(t, usernames, []string{"foo@email", "bar@email"})
}

func TestSaml_GetUsernameFromEmail(t *testing.T) { //nolint:dupl
	t.Parallel()

	ctx := context.TODO()

	gitHubClient := mocks.NewIGitHubV4Saml(t)
	discovery := New(nil, "test")
	discovery.client = gitHubClient

	queryFn := func(name string) func(ctx context.Context, q interface{}, variables map[string]interface{}) {
		return func(ctx context.Context, query interface{}, variables map[string]interface{}) {
			resp := fmt.Sprintf(`
{
  "Organization": {
    "SamlIdentityProvider": {
      "ExternalIdentities": {
        "Edges": [
          {
            "Node": {
              "User": {
                "Login": "%s"
              }
            }
          }
        ]
      }
    }
  }
}
`, name)
			arg := query.(*usernameQuery)
			err := json.Unmarshal([]byte(resp), arg)

			assert.NoError(t, err)
			assert.Equal(t, name,
				arg.Organization.SamlIdentityProvider.ExternalIdentities.Edges[0].Node.User.Login,
			)
		}
	}

	gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
		"email": githubv4.String("foo@email"), "org": githubv4.String("test")},
	).Run(queryFn("foo")).Return(nil)
	gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
		"email": githubv4.String("bar@email"), "org": githubv4.String("test")},
	).Run(queryFn("bar")).
		Return(nil)

	usernames, err := discovery.GetUsernameFromEmail(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
	assert.ElementsMatch(t, usernames, []string{"foo", "bar"})
}
