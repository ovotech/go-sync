package team_test

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/terraformcloud/team"
)

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

func ExampleWithClient() {
	ctx := context.Background()

	client, err := tfe.NewClient(&tfe.Config{Token: "token"})
	if err != nil {
		log.Fatal(err)
	}

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.Organisation: "ovotech",
	}, team.WithClient(client))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.Token:        "my-org-token",
		team.Organisation: "ovotech",
	}, team.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
