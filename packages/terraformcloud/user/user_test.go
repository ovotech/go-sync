package user_test

import (
	"context"
	"log"

	"github.com/hashicorp/go-tfe"

	"github.com/ovotech/go-sync/packages/gosync"
	"github.com/ovotech/go-sync/packages/terraformcloud/user"
)

func ExampleNew() {
	client, err := tfe.NewClient(&tfe.Config{Token: "my-org-token"})
	if err != nil {
		log.Fatal(err)
	}

	adapter := user.New(client, "my-org", "my-team")

	gosync.New(adapter)
}

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
