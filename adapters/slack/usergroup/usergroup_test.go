package usergroup_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/slack/usergroup"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := usergroup.Init(ctx, map[gosync.ConfigKey]string{
		usergroup.SlackAPIKey: "my-slack-token",
		usergroup.UserGroupID: "S0123ABC456",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
