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
	"math"
	"os"
	"strings"
	"time"

	"github.com/slack-go/slack"

	gosync "github.com/ovotech/go-sync"
)

/*
SlackAPIKey is an API key for authenticating with Slack.
*/
const SlackAPIKey gosync.ConfigKey = "slack_api_key" //nolint:gosec

// Name is the Slack conversation name.
const Name gosync.ConfigKey = "name"

/*
MuteRestrictedErrOnKickFromPublic mutes an error that occurs when Slack is configured to prevent kicking users from
public conversations. Set this to true to mute this error and continue syncing.
*/
const MuteRestrictedErrOnKickFromPublic gosync.ConfigKey = "mute_restricted_err_kick_from_public"

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
	MuteRestrictedErrOnKickFromPublic bool // See [conversation.MuteRestrictedErrOnKickFromPublic]
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

// paginateUsersInfo requests.
func (c *Conversation) paginateUsersInfo(slackUsers []string) (*[]slack.User, error) {
	currentPage := 0
	pageSize := 30
	out := make([]slack.User, 0, cap(slackUsers))
	totalPages := math.Floor(float64(len(slackUsers)) / float64(pageSize))

	for {
		c.Logger.Printf("Calling GetUsersInfo page %v of %v", currentPage+1, totalPages+1)

		start := currentPage * pageSize
		end := (currentPage * pageSize) + pageSize

		if end > cap(slackUsers) {
			end = cap(slackUsers)
		}

		// Get a page of slackUsers to send up to the API.
		page := slackUsers[start:end]

		// Request only the page of users.
		users, err := c.client.GetUsersInfo(page...)
		if err != nil {
			return nil, fmt.Errorf("paginateusersinfo -> %w", err)
		}

		// Append the results to combined output.
		out = append(out, *users...)

		// When the output matches the number of input slack users, end the loop.
		if len(out) == len(slackUsers) {
			break
		}

		// Increment the current page.
		currentPage++
	}

	return &out, nil
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

	if len(slackUsers) == 0 {
		c.Logger.Println("Fetched no accounts from conversation")

		return []string{}, nil
	}

	users, err := c.paginateUsersInfo(slackUsers)
	if err != nil {
		return nil, fmt.Errorf("slack.conversation.get.getusersinfo -> %w", err)
	}

	emails := make([]string, 0, len(*users))

	for _, user := range *users {
		if !user.IsBot {
			emails = append(emails, strings.ToLower(user.Profile.Email))

			// Add the email -> ID map for use with Remove method.
			c.cache[strings.ToLower(user.Profile.Email)] = user.ID
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
		user, err := c.client.GetUserByEmail(strings.ToLower(email))
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
		err := c.client.KickUserFromConversation(c.conversationName, c.cache[strings.ToLower(email)])
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

// WithClient passes a custom Slack client to the adapter.
func WithClient(client *slack.Client) gosync.ConfigFn {
	return func(i interface{}) {
		if adapter, ok := i.(*Conversation); ok {
			adapter.client = client
		}
	}
}

// WithLogger passes a custom logger to the adapter.
func WithLogger(logger *log.Logger) gosync.ConfigFn {
	return func(i interface{}) {
		if adapter, ok := i.(*Conversation); ok {
			adapter.Logger = logger
		}
	}
}

/*
Init a new Slack Conversation [gosync.Adapter].

Required config:
  - [conversation.Name]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string, configFns ...gosync.ConfigFn) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{Name} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("slack.conversation.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	adapter := &Conversation{
		MuteRestrictedErrOnKickFromPublic: false,
		conversationName:                  config[Name],
		cache:                             make(map[string]string),
	}

	if _, ok := config[SlackAPIKey]; ok {
		client := slack.New(config[SlackAPIKey])

		WithClient(client)(adapter)
	}

	for _, configFn := range configFns {
		configFn(adapter)
	}

	if val, ok := config[MuteRestrictedErrOnKickFromPublic]; ok {
		adapter.MuteRestrictedErrOnKickFromPublic = strings.ToLower(val) == "true"
	}

	if adapter.Logger == nil {
		logger := log.New(
			os.Stderr, "[go-sync/slack/conversation] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		)

		WithLogger(logger)(adapter)
	}

	if adapter.client == nil {
		return nil, fmt.Errorf("slack.conversation.init -> %w(%s)", gosync.ErrMissingConfig, SlackAPIKey)
	}

	return adapter, nil
}
