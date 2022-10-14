# GitHub Team adapter for Go Sync
This adapter synchronises email addresses with a GitHub team.

## Requirements
In order to synchronise with GitHub, you'll need to create a [Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
with the following permissions:

| Scopes                            |
|-----------------------------------|
| admin:org                         |
| write:org                         |
| read:org                          |

## Example
```go
package main

import (
	"context"
	"log"

	"github.com/google/go-github/v47/github"
	"github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/github/discovery/saml"
	"github.com/ovotech/go-sync/adapters/github/team"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	// Authenticated client to communicate with GitHub APIs.
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "my-github-token"},
	))

	var (
		gitHubV3Client = github.NewClient(oauthClient)      // GitHub V3 API is used by GH Teams adapter.
		gitHubV4Client = githubv4.NewClient(oauthClient)    // GitHub V4 API is used by SAML discovery.
		samlClient     = saml.New(gitHubV4Client, "my-org") // GitHub Discovery service uses SAML to convert emails into GH users.
	)

	ghTeam := team.New(gitHubV3Client, samlClient, "my-org", "my-team-slug")

	svc := gosync.New(ghTeam)

	// Synchronise a Slack User group with something else.
	anotherServiceAdapter := someAdapter.New()

	err := svc.SyncWith(context.Background(), anotherServiceAdapter)
	if err != nil {
		log.Fatal(err)
	}
}
```
