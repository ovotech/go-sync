package usergroup_test

import (
	"context"
	"log"
	"os"

	"github.com/slack-go/slack"

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

func ExampleWithClient() {
	ctx := context.Background()

	client := slack.New("my-slack-token")

	adapter, err := usergroup.Init(ctx, map[gosync.ConfigKey]string{
		usergroup.UserGroupID: "S0123ABC456",
	}, usergroup.WithClient(client))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := usergroup.Init(ctx, map[gosync.ConfigKey]string{
		usergroup.SlackAPIKey: "my-slack-token",
		usergroup.UserGroupID: "S0123ABC456",
	}, usergroup.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
