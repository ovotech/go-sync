package user_test

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/terraformcloud/user"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := user.Init(ctx, map[gosync.ConfigKey]string{
		user.Token:        "my-org-token",
		user.Organisation: "ovotech",
		user.Team:         "my-team",
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

	adapter, err := user.Init(ctx, map[gosync.ConfigKey]string{
		user.Organisation: "ovotech",
		user.Team:         "my-team",
	}, user.WithClient(client))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := user.Init(ctx, map[gosync.ConfigKey]string{
		user.Token:        "my-org-token",
		user.Organisation: "ovotech",
		user.Team:         "my-team",
	}, user.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
