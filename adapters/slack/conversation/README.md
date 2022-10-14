# Slack Conversation adapter for Go Sync
This adapter synchronises email addresses with a Slack conversation.

## Warning
The Slack usergroup API doesn't allow a usergroup to have no members. If this behaviour is expected, we recommend
setting `adapter.MuteGroupCannotBeEmpty = true` to mute the error. No members will be removed, but Go Sync will continue
processing.

## Requirements
In order to synchronise with Slack, you'll need to [create a Slack app](https://api.slack.com/authentication/basics)
with the following OAuth permissions:

| Bot Token Scopes                                                  |
|-------------------------------------------------------------------|
| [users:read](https://api.slack.com/scopes/users:read)             |
| [users:read.email](https://api.slack.com/scopes/users:read.email) |
| [channels:manage](https://api.slack.com/scopes/channels:manage)   |
| [channels:read](https://api.slack.com/scopes/channels:read)       |
| [groups:read](https://api.slack.com/scopes/groups:read)           |
| [groups:write](https://api.slack.com/scopes/groups:write)         |
| [im:write](https://api.slack.com/scopes/im:write)                 |
| [mpim:write](https://api.slack.com/scopes/mpim:write)             |

## Example
```go
package main

import (
	"context"
	"log"

	"github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/slack/conversation"
	"github.com/slack-go/slack"
)

func main() {
	slackClient := slack.New("my-slack-token")
	conversationAdapter := conversation.New(slackClient, "UG000123")
	
	svc := gosync.New(conversationAdapter)

	// Synchronise a Slack User group with something else.
	anotherServiceAdapter := someAdapter.New()

	err := svc.SyncWith(context.Background(), anotherServiceAdapter)
	if err != nil {
		log.Fatal(err)
	}
}
```
