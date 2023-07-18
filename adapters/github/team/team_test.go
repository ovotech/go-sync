package team_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/github/team"
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
