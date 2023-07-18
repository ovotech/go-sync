package conversation_test

import (
	"context"
	"log"

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
