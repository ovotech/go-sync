package team_test

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v47/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/ovotech/go-sync/packages/github/discovery/saml"
	"github.com/ovotech/go-sync/packages/github/team"
	"github.com/ovotech/go-sync/packages/gosync"
)

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

func ExampleWithGitHubV3Client() {
	ctx := context.Background()

	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "my token"},
	))

	client := github.NewClient(oauthClient)

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.GitHubOrg:          "my-org",
		team.TeamSlug:           "my-team-slug",
		team.DiscoveryMechanism: "saml",
	}, team.WithGitHubV3Client(client))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithDiscoveryService() {
	ctx := context.Background()

	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "my token"},
	))

	client := githubv4.NewClient(oauthClient)

	samlSvc := saml.New(client, "my-org")

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.GitHubToken: "my-github-token",
		team.GitHubOrg:   "my-org",
		team.TeamSlug:    "my-team-slug",
	}, team.WithDiscoveryService(samlSvc))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stderr, "my-custom-logger ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.GitHubToken: "my-github-token",
		team.GitHubOrg:   "my-org",
		team.TeamSlug:    "my-team-slug",
	}, team.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
