package user_test

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"

	"github.com/ovotech/go-sync/adapters/terraformcloud/user"
	"github.com/ovotech/go-sync/internal/gosync"
	"github.com/ovotech/go-sync/pkg/types"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := user.Init(ctx, map[types.ConfigKey]string{
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

	adapter, err := user.Init(ctx, map[types.ConfigKey]string{
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

	adapter, err := user.Init(ctx, map[types.ConfigKey]string{
		user.Token:        "my-org-token",
		user.Organisation: "ovotech",
		user.Team:         "my-team",
	}, user.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
