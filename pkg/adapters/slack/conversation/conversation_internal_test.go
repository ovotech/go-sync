package conversation

import (
	"context"
	"testing"

	"github.com/ovotech/go-sync/mocks"
	"github.com/ovotech/go-sync/pkg/ports"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestImplementsInterface(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*ports.Adapter)(nil), &Conversation{})
}

func TestNew(t *testing.T) {
	t.Parallel()

	slackClient := mocks.NewISlackConversation(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	assert.Equal(t, "test", adapter.conversationName)
	assert.Zero(t, slackClient.Calls)
}

func TestConversation_Get(t *testing.T) {
	t.Parallel()

	slackClient := mocks.NewISlackConversation(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	// First page.
	slackClient.EXPECT().GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: "test",
		Cursor:    "",
		Limit:     50,
	}).Return([]string{"slack-foo"}, "page-2", nil)

	// Second page.
	slackClient.EXPECT().GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: "test",
		Cursor:    "page-2",
		Limit:     50,
	}).Return([]string{"slack-bar"}, "", nil)

	// Users info response.
	slackClient.EXPECT().GetUsersInfo("slack-foo", "slack-bar").Return(&[]slack.User{
		{ID: "foo", IsBot: false, Profile: slack.UserProfile{Email: "foo@email"}},
		{ID: "bar", IsBot: false, Profile: slack.UserProfile{Email: "bar@email"}},
	}, nil)

	accounts, err := adapter.Get(context.TODO())

	assert.NoError(t, err)
	assert.ElementsMatch(t, accounts, []string{"foo@email", "bar@email"})
	assert.Equal(t, map[string]string{"foo@email": "foo", "bar@email": "bar"}, adapter.cache)
}

func TestConversation_Add(t *testing.T) {
	t.Parallel()

	slackClient := mocks.NewISlackConversation(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	slackClient.EXPECT().GetUserByEmail("foo@email").Return(&slack.User{
		ID: "foo",
	}, nil)
	slackClient.EXPECT().GetUserByEmail("bar@email").Return(&slack.User{
		ID: "bar",
	}, nil)
	slackClient.EXPECT().InviteUsersToConversation("test", "foo", "bar").Return(nil, nil)

	err := adapter.Add(context.TODO(), []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"foo@email": "foo", "bar@email": "bar"}, adapter.cache)
}

func TestConversation_Remove(t *testing.T) {
	t.Parallel()

	slackClient := mocks.NewISlackConversation(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient
	adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

	slackClient.EXPECT().KickUserFromConversation("test", "foo").Return(nil)
	slackClient.EXPECT().KickUserFromConversation("test", "bar").Return(nil)

	err := adapter.Remove(context.TODO(), []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{}, adapter.cache)
}
