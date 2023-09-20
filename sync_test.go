package gosync_test

import (
	"context"
	"log"
	"os"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/github/team"
	"github.com/ovotech/go-sync/adapters/slack/conversation"
)

func ExampleNew() {
	ctx := context.Background()

	// Specify a custom logger for the GitHub adapter.
	logger := log.New(os.Stderr, "new logger", log.LstdFlags)

	// Create a GitHub team adapter.
	source, err := team.Init(ctx, map[gosync.ConfigKey]string{
		team.GitHubToken: "some-token",
	}, team.WithLogger(logger))
	if err != nil {
		log.Panic(err)
	}

	// Create a Slack conversation adapter.
	destination, err := conversation.Init(ctx, map[gosync.ConfigKey]string{
		conversation.Name: "example",
	})
	if err != nil {
		log.Panic(err)
	}

	err = gosync.New(source).SyncWith(ctx, destination)
	if err != nil {
		log.Panic(err)
	}
}
