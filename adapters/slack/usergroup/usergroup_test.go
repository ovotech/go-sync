package usergroup_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/slack/usergroup"
	"github.com/slack-go/slack"
)

func ExampleNew() {
	client := slack.New("my-slack-token")

	adapter := usergroup.New(client, "S0123ABC456")

	gosync.New(adapter)
}

func ExampleInit() {
	ctx := context.Background()

	adapter, err := usergroup.Init(ctx, map[gosync.ConfigKey]string{
		usergroup.SlackAPIKey:      "my-slack-token",
		usergroup.SlackUserGroupID: "S0123ABC456",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
