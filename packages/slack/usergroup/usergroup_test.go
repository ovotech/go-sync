package usergroup_test

import (
	"context"
	"log"

	"github.com/slack-go/slack"

	"github.com/ovotech/go-sync/packages/gosync"
	"github.com/ovotech/go-sync/packages/slack/usergroup"
)

func ExampleNew() {
	client := slack.New("my-slack-token")

	adapter := usergroup.New(client, "S0123ABC456")

	gosync.New(adapter)
}

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
