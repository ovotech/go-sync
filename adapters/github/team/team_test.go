package team_test

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v47/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/ovotech/go-sync/adapters/github/discovery/saml"
	"github.com/ovotech/go-sync/adapters/github/team"
	"github.com/ovotech/go-sync/internal/gosync"
	"github.com/ovotech/go-sync/pkg/types"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := team.Init(ctx, map[types.ConfigKey]string{
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

func ExampleWithClient() {
	ctx := context.Background()

	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "token"},
	))

	gitHubClient := github.NewClient(oauthClient)

	adapter, err := team.Init(ctx, map[types.ConfigKey]string{
		team.GitHubToken:        "my-github-token",
		team.GitHubOrg:          "my-org",
		team.TeamSlug:           "my-team-slug",
		team.DiscoveryMechanism: "saml",
	}, team.WithClient(gitHubClient))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithDiscoveryService() {
	ctx := context.Background()

	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "token"},
	))

	discoverySvc := saml.New(githubv4.NewClient(oauthClient), "my-org")

	adapter, err := team.Init(ctx, map[types.ConfigKey]string{
		team.GitHubToken: "my-github-token",
		team.GitHubOrg:   "my-org",
		team.TeamSlug:    "my-team-slug",
	}, team.WithDiscoveryService(discoverySvc))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := team.Init(ctx, map[types.ConfigKey]string{
		team.GitHubToken:        "my-github-token",
		team.GitHubOrg:          "my-org",
		team.TeamSlug:           "my-team-slug",
		team.DiscoveryMechanism: "saml",
	}, team.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
