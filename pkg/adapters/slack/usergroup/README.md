# Slack UserGroup adapter for Go Sync
This adapter synchronises email addresses with a Slack User group.

## Example
```go

package main

import (
	"context"
	"log"

	"github.com/ovotech/go-sync/pkg/adapters/slack/usergroup"
	"github.com/ovotech/go-sync/pkg/sync"
	"github.com/slack-go/slack"
)

func main() {
	slackClient := slack.New("my-slack-token")
	userGroupAdapter := usergroup.New(slackClient, "UG000123")
	
	svc := sync.New(userGroupAdapter)

	// Synchronise a Slack User group with something else.
	anotherServiceAdapter := someAdapter.New()

	err := svc.SyncWith(context.Background(), anotherServiceAdapter)
	if err != nil {
		log.Fatal(err)
	}
}
```

[Information on how to obtain a Slack token.](../README.md#requirements)