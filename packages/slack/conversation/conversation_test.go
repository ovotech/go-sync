package conversation_test

import (
	"context"
	"log"

	"github.com/slack-go/slack"

	"github.com/ovotech/go-sync/packages/gosync"
	"github.com/ovotech/go-sync/packages/slack/conversation"
)

func ExampleNew() {
	client := slack.New("my-slack-token")

	adapter := conversation.New(client, "C0123ABC456")

	gosync.New(adapter)
}

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
