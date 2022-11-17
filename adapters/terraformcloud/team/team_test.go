package team_test

import (
	"context"
	"log"

	"github.com/hashicorp/go-tfe"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/terraformcloud/team"
)

func ExampleNew() {
	client, err := tfe.NewClient(&tfe.Config{Token: "my-org-token"})
	if err != nil {
		log.Fatal(err)
	}

	adapter := team.New(client, "my-org")

	gosync.New(adapter)
}

func ExampleInit() {
	ctx := context.Background()

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.Token:        "my-org-token",
		team.Organisation: "ovotech",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
