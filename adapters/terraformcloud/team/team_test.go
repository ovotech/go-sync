package team_test

import (
	"context"
	"log"

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
