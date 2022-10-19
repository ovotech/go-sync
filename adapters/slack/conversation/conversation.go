/*
Package conversation synchronises email addresses with a Slack conversation.

# Requirements

In order to synchronise with Slack, you'll need to [create a Slack app] with the following OAuth Bot Token permissions:
  - [users:read]
  - [users:read.email]
  - [channels:manage]
  - [channels:read]
  - [groups:read]
  - [groups:write]
  - [im:write]
  - [mpim:write]

# Examples

See [New] and [Init].

[create a Slack app]: https://api.slack.com/authentication/basics
[users:read]: https://api.slack.com/scopes/users:read
[users:read.email]: https://api.slack.com/scopes/users:read.email
[channels:manage]: https://api.slack.com/scopes/channels:manage
[channels:read]: https://api.slack.com/scopes/channels:read
[groups:read]: https://api.slack.com/scopes/groups:read
[groups:write]: https://api.slack.com/scopes/groups:write
[im:write]: https://api.slack.com/scopes/im:write
[mpim:write]: https://api.slack.com/scopes/mpim:write
*/
package conversation

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gosync "github.com/ovotech/go-sync"
	"github.com/slack-go/slack"
)

/*
SlackAPIKey is an API key for authenticating with Slack.
*/
const SlackAPIKey gosync.ConfigKey = "slack_api_key" //nolint:gosec

// SlackConversationName is the Slack conversation name.
const SlackConversationName gosync.ConfigKey = "slack_conversation_name"

var (
	// Ensure [conversation.Conversation] fully satisfies the [gosync.Adapter] interface.
	_ gosync.Adapter = &Conversation{}
	_ gosync.InitFn  = Init // Ensure the [conversation.Init] function fully satisfies the [gosync.InitFn] type.
)

// iSlackConversation is a subset of the Slack Client, and used to build mocks for easy testing.
type iSlackConversation interface {
	GetUsersInConversation(params *slack.GetUsersInConversationParameters) ([]string, string, error)
	GetUsersInfo(users ...string) (*[]slack.User, error)
	GetUserByEmail(email string) (*slack.User, error)
	InviteUsersToConversation(channelID string, users ...string) (*slack.Channel, error)
	KickUserFromConversation(channelID string, user string) error
}

type Conversation struct {
	/*
		Slack may be configured to only allow admins to kick from public conversations, which will fail the entire sync
		job. Set to true to mute this error and continue synchronisation.
	*/
	MuteRestrictedErrOnKickFromPublic bool
	client                            iSlackConversation
	conversationName                  string
	// cache stores the Slack ID -> email mapping for use with the Remove method.
	cache  map[string]string
	Logger *log.Logger
}

// getListOfSlackUsernames gets a list of Slack users in a conversation, and paginates through the results.
func (c *Conversation) getListOfSlackUsernames() ([]string, error) {
	var (
		cursor string
		users  []string
		err    error
	)

	for {
		params := &slack.GetUsersInConversationParameters{
			ChannelID: c.conversationName,
			Cursor:    cursor,
			Limit:     50, //nolint:gomnd
		}

		var pageOfUsers []string

		pageOfUsers, cursor, err = c.client.GetUsersInConversation(params)
		if err != nil {
			return nil, fmt.Errorf("getusersinconversation(%s) -> %w", c.conversationName, err)
		}

		users = append(users, pageOfUsers...)

		if cursor == "" {
			break
		}
	}

	return users, nil
}

// Get email addresses in a Slack Conversation.
func (c *Conversation) Get(_ context.Context) ([]string, error) {
	c.Logger.Printf("Fetching accounts from Slack conversation %s", c.conversationName)

	// Initialise the cache.
	c.cache = make(map[string]string)

	slackUsers, err := c.getListOfSlackUsernames()
	if err != nil {
		return nil, fmt.Errorf("slack.conversation.get.getlistofslackusernames -> %w", err)
	}

	users, err := c.client.GetUsersInfo(slackUsers...)
	if err != nil {
		return nil, fmt.Errorf("slack.conversation.get.getusersinfo -> %w", err)
	}

	emails := make([]string, 0, len(*users))

	for _, user := range *users {
		if !user.IsBot {
			emails = append(emails, user.Profile.Email)

			// Add the email -> ID map for use with Remove method.
			c.cache[user.Profile.Email] = user.ID
		}
	}

	c.Logger.Println("Fetched accounts successfully")

	return emails, nil
}

// Add email addresses to a Slack Conversation.
func (c *Conversation) Add(_ context.Context, emails []string) error {
	c.Logger.Printf("Adding %s to Slack conversation %s", emails, c.conversationName)

	slackIds := make([]string, len(emails))

	for index, email := range emails {
		user, err := c.client.GetUserByEmail(email)
		if err != nil {
			return fmt.Errorf("slack.conversation.add.getuserbyemail(%s) -> %w", email, err)
		}

		slackIds[index] = user.ID
	}

	_, err := c.client.InviteUsersToConversation(c.conversationName, slackIds...)
	if err != nil {
		return fmt.Errorf("slack.conversation.add.inviteuserstoconversation(%s, ...) -> %w", c.conversationName, err)
	}

	c.Logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove email addresses from a Slack Conversation.
func (c *Conversation) Remove(_ context.Context, emails []string) error {
	c.Logger.Printf("Removing %s from Slack conversation %s", emails, c.conversationName)

	// If the cache hasn't been generated, regenerate it.
	if c.cache == nil {
		return fmt.Errorf("slack.conversation.remove -> %w", gosync.ErrCacheEmpty)
	}

	for _, email := range emails {
		err := c.client.KickUserFromConversation(c.conversationName, c.cache[email])
		if err != nil {
			if c.MuteRestrictedErrOnKickFromPublic && strings.Contains(err.Error(), "restricted_action") {
				c.Logger.Println("Cannot kick from public channel, but error is muted by configuration - continuing")

				return nil
			}

			return fmt.Errorf(
				"slack.conversation.remove.kickuserfromconversation(%s, %s) -> %w",
				c.conversationName,
				c.cache[email],
				err,
			)
		}

		// To prevent rate limiting, sleep for 1 second after each kick.
		time.Sleep(1 * time.Second)
	}

	c.Logger.Println("Finished removing accounts successfully")

	return nil
}

// New Slack Conversation [gosync.Adapter].
func New(client *slack.Client, conversationName string, optsFn ...func(conversation *Conversation)) *Conversation {
	conversation := &Conversation{
		MuteRestrictedErrOnKickFromPublic: false,
		client:                            client,
		conversationName:                  conversationName,
		cache:                             nil,
		Logger: log.New(
			os.Stderr,
			"[go-sync/slack/conversation] ",
			log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		),
	}

	for _, fn := range optsFn {
		fn(conversation)
	}

	return conversation
}

/*
Init a new Slack Conversation [gosync.Adapter].

Required config:
  - [conversation.SlackAPIKey]
  - [conversation.SlackConversationName]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{SlackAPIKey, SlackConversationName} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("slack.conversation.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	client := slack.New(config[SlackAPIKey])

	return New(client, config[SlackConversationName]), nil
}
