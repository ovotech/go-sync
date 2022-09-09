package slack

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

type Conversation struct {
	client           *slack.Client
	conversationName string
	cache            map[string]string // This stores the Slack ID -> email mapping for use with the Remove method.
}

// NewConversationService instantiates a new Slack conversation service.
func NewConversationService(client *slack.Client, channelName string) *Conversation {
	return &Conversation{
		client:           client,
		conversationName: channelName,
		cache:            map[string]string{},
	}
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
			return nil, fmt.Errorf("GetUsersInConversation(%s) -> %w", c.conversationName, err)
		}

		users = append(users, pageOfUsers...)

		if cursor == "" {
			break
		}
	}

	return users, nil
}

func (c *Conversation) get() ([]string, error) {
	slackUsers, err := c.getListOfSlackUsernames()
	if err != nil {
		return nil, fmt.Errorf("getListOfSlackUsernames -> %w", err)
	}

	users, err := c.client.GetUsersInfo(slackUsers...)
	if err != nil {
		return nil, fmt.Errorf("GetUsersInfo -> %w", err)
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

// Get gets a list of emails from a Slack channel.
func (c *Conversation) Get() ([]string, error) {
	users, err := c.get()
	if err != nil {
		return nil, fmt.Errorf("slack.conversation.Get -> %w", err)
	}

	return users, nil
}

// Add adds an email to a conversation.
func (c *Conversation) Add(emails ...string) ([]string, []error, error) {
	slackIds := make([]string, len(emails))

	for index, email := range emails {
		user, err := c.client.GetUserByEmail(email)
		if err != nil {
			return nil, nil, fmt.Errorf("slack.conversation.Add.GetUserByEmail(%s) -> %w", email, err)
		}

		slackIds[index] = user.ID
		// Add the user to the cache.
		c.cache[email] = user.ID
	}

	if _, err := c.client.InviteUsersToConversation(c.conversationName, slackIds...); err != nil {
		c.cache = map[string]string{}

		return nil, nil, fmt.Errorf(
			"slack.conversation.Add.InviteUsersToConversation(%s, ...) -> %w",
			c.conversationName,
			err,
		)
	}

	return emails, nil, nil
}

// Remove removes email addresses from a conversation.
func (c *Conversation) Remove(emails ...string) ([]string, []error, error) {
	var (
		success []string
		failure []error
	)

	// If the cache hasn't been generated, regenerate it.
	if len(c.cache) == 0 {
		if _, err := c.get(); err != nil {
			//goland:noinspection ALL
			return nil, nil, fmt.Errorf("slack.conversation.Remove.get -> %w", err)
		}
	}

	for _, email := range emails {
		err := c.client.KickUserFromConversation(c.conversationName, c.cache[email])
		if err != nil {
			failure = append(
				failure,
				fmt.Errorf(
					"slack.conversation.Remove.KickUserFromConversation(%s, %s) -> %w",
					c.conversationName,
					c.cache[email],
					err,
				),
			)
		} else {
			success = append(success, email)
			// Delete the entry from the cache.
			delete(c.cache, email)
		}

		// To prevent rate limiting, sleep for 1 second after each kick.
		time.Sleep(1 * time.Second)
	}

	return success, failure, nil
}
