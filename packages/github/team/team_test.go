package team_test

import (
	"context"
	"log"

	"github.com/google/go-github/v47/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/ovotech/go-sync/packages/github/discovery/saml"
	"github.com/ovotech/go-sync/packages/github/team"
	"github.com/ovotech/go-sync/packages/gosync"
)

func ExampleNew() {
	ctx := context.Background()

	// Authenticated client to communicate with GitHub APIs.
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "my-github-token"},
	))

	var (
		// GitHub V3 API is used by GH Teams adapter.
		gitHubV3Client = github.NewClient(oauthClient)
		// GitHub V4 API is used by SAML discovery.
		gitHubV4Client = githubv4.NewClient(oauthClient)
		// GitHub Discovery service uses SAML to convert emails into GH users.
		samlClient = saml.New(gitHubV4Client, "my-org")
	)

	adapter := team.New(gitHubV3Client, samlClient, "my-org", "my-team-slug")

	gosync.New(adapter)
}

func ExampleInit() {
	ctx := context.Background()

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.GitHubToken:        "my-github-token",
		team.GitHubOrg:          "my-org",
		team.TeamSlug:           "my-team-slug",
		team.DiscoveryMechanism: "saml",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
