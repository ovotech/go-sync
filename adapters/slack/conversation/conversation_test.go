package conversation_test

import (
	"context"
	"log"
	"os"

	"github.com/slack-go/slack"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/slack/conversation"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := conversation.Init(ctx, map[gosync.ConfigKey]string{
		conversation.SlackAPIKey: "my-slack-token",
		conversation.Name:        "C0123ABC456",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithClient() {
	ctx := context.Background()

	client := slack.New("my-slack-token")

	adapter, err := conversation.Init(ctx, map[gosync.ConfigKey]string{
		conversation.Name: "C0123ABC456",
	}, conversation.WithClient(client))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := conversation.Init(ctx, map[gosync.ConfigKey]string{
		conversation.SlackAPIKey: "my-slack-token",
		conversation.Name:        "C0123ABC456",
	}, conversation.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
