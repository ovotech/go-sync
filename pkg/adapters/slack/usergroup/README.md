# Slack UserGroup adapter for Go Sync
This adapter synchronises email addresses with a Slack User group.

## Requirements
In order to synchronise with Slack, you'll need to [create a Slack app](https://api.slack.com/authentication/basics)
with the following OAuth permissions:

| Bot Token Scopes                                                  |
|-------------------------------------------------------------------|
| [users:read](https://api.slack.com/scopes/users:read)             |
| [users:read.email](https://api.slack.com/scopes/users:read.email) |
| [usergroups:read](https://api.slack.com/scopes/usergroups:read)   |
| [usergroups:write](https://api.slack.com/scopes/usergroups:write) |

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
