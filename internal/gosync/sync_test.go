package gosync_test

import (
	"context"
	"log"
	"os"

	"github.com/ovotech/go-sync/adapters/github/team"
	"github.com/ovotech/go-sync/adapters/slack/conversation"
	"github.com/ovotech/go-sync/internal/gosync"
	"github.com/ovotech/go-sync/pkg/types"
)

func ExampleNew() {
	ctx := context.Background()

	// Specify a custom logger for the GitHub adapter.
	logger := log.New(os.Stderr, "new logger", log.LstdFlags)

	// Create a GitHub team adapter.
	source, err := team.Init(ctx, map[types.ConfigKey]string{
		team.GitHubToken: "some-token",
	}, team.WithLogger(logger))
	if err != nil {
		log.Panic(err)
	}

	// Create a Slack conversation adapter.
	destination, err := conversation.Init(ctx, map[types.ConfigKey]string{
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
