package conversation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

type iSlackConversation interface {
	GetUsersInConversation(params *slack.GetUsersInConversationParameters) ([]string, string, error)
	GetUsersInfo(users ...string) (*[]slack.User, error)
	GetUserByEmail(email string) (*slack.User, error)
	InviteUsersToConversation(channelID string, users ...string) (*slack.Channel, error)
	KickUserFromConversation(channelID string, user string) error
}

type Conversation struct {
	client           iSlackConversation
	conversationName string
	cache            map[string]string // This stores the Slack ID -> email mapping for use with the Remove method.
}

var ErrCacheEmpty = errors.New("cache is empty - run Get()")

// New instantiates a new Slack conversation adapter.
func New(client *slack.Client, channelName string, optsFn ...func(conversation *Conversation)) *Conversation {
	conversation := &Conversation{
		client:           client,
		conversationName: channelName,
		cache:            make(map[string]string),
	}

	for _, fn := range optsFn {
		fn(conversation)
	}

	return conversation
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

// Get gets a list of emails from a Slack channel.
func (c *Conversation) Get(_ context.Context) ([]string, error) {
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

	return emails, nil
}

// Add adds an email to a conversation.
func (c *Conversation) Add(_ context.Context, emails []string) error {
	slackIds := make([]string, len(emails))

	for index, email := range emails {
		user, err := c.client.GetUserByEmail(email)
		if err != nil {
			return fmt.Errorf("slack.conversation.add.getuserbyemail(%s) -> %w", email, err)
		}

		slackIds[index] = user.ID
		// Add the user to the cache.
		c.cache[email] = user.ID
	}

	_, err := c.client.InviteUsersToConversation(c.conversationName, slackIds...)
	if err != nil {
		c.cache = nil

		return fmt.Errorf("slack.conversation.add.inviteuserstoconversation(%s, ...) -> %w", c.conversationName, err)
	}

	return nil
}

// Remove removes email addresses from a conversation.
func (c *Conversation) Remove(_ context.Context, emails ...string) error {
	// If the cache hasn't been generated, regenerate it.
	if len(c.cache) == 0 {
		return fmt.Errorf("slack.conversation.remove -> %w", ErrCacheEmpty)
	}

	for _, email := range emails {
		err := c.client.KickUserFromConversation(c.conversationName, c.cache[email])
		if err != nil {
			return fmt.Errorf(
				"slack.conversation.remove.kickuserfromconversation(%s, %s) -> %w",
				c.conversationName,
				c.cache[email],
				err,
			)
		}

		// Delete the entry from the cache.
		delete(c.cache, email)

		// To prevent rate limiting, sleep for 1 second after each kick.
		time.Sleep(1 * time.Second)
	}

	return nil
}
