/*
Package conversation synchronises email addresses with Slack conversations.

In order to use this adapter, you'll need an authenticated Slack client and for the Slack app to have been added
to the conversation.
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

// Ensure the Slack Conversation adapter type fully satisfies the gosync.Adapter interface.
var _ gosync.Adapter = &Conversation{}

const (
	/*
		API key for authenticating with Slack.
	*/
	SlackAPIKey           gosync.ConfigKey = "slack_api_key"           //nolint:gosec
	SlackConversationName gosync.ConfigKey = "slack_conversation_name" // Slack conversation name.
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
	// Slack may be configured to only allow admins to kick from public conversations, which will fail the entire sync
	// job. Set to true to mute this error and continue synchronisation.
	MuteRestrictedErrOnKickFromPublic bool
	client                            iSlackConversation
	conversationName                  string
	// cache stores the Slack ID -> email mapping for use with the Remove method.
	cache  map[string]string
	logger *log.Logger
}

// WithLogger sets a custom logger.
func WithLogger(logger *log.Logger) func(*Conversation) {
	return func(conversation *Conversation) {
		conversation.logger = logger
	}
}

// New instantiates a new Slack Conversation adapter.
func New(client *slack.Client, conversationName string, optsFn ...func(conversation *Conversation)) *Conversation {
	conversation := &Conversation{
		MuteRestrictedErrOnKickFromPublic: false,
		client:                            client,
		conversationName:                  conversationName,
		cache:                             nil,
		logger: log.New(
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

// Ensure the Init function fully satisfies the gosync.InitFn type.
var _ gosync.InitFn = Init

// Init a new Slack Conversation gosync.Adapter. All gosync.ConfigKey keys are required in config.
func Init(config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{SlackAPIKey, SlackConversationName} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("slack.conversation.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	client := slack.New(config[SlackAPIKey])

	return New(client, config[SlackConversationName]), nil
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

// Get emails of Slack users in a conversation.
func (c *Conversation) Get(_ context.Context) ([]string, error) {
	c.logger.Printf("Fetching accounts from Slack conversation %s", c.conversationName)

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

	c.logger.Println("Fetched accounts successfully")

	return emails, nil
}

// Add emails to a Slack conversation.
func (c *Conversation) Add(_ context.Context, emails []string) error {
	c.logger.Printf("Adding %s to Slack conversation %s", emails, c.conversationName)

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

	c.logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove emails from a Slack conversation.
func (c *Conversation) Remove(_ context.Context, emails []string) error {
	c.logger.Printf("Removing %s from Slack conversation %s", emails, c.conversationName)

	// If the cache hasn't been generated, regenerate it.
	if c.cache == nil {
		return fmt.Errorf("slack.conversation.remove -> %w", gosync.ErrCacheEmpty)
	}

	for _, email := range emails {
		err := c.client.KickUserFromConversation(c.conversationName, c.cache[email])
		if err != nil {
			if c.MuteRestrictedErrOnKickFromPublic && strings.Contains(err.Error(), "restricted_action") {
				c.logger.Println("Cannot kick from public channel, but error is muted by configuration - continuing")

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

	c.logger.Println("Finished removing accounts successfully")

	return nil
}
