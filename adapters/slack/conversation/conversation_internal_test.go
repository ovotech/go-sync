package conversation

import (
	"context"
	"errors"
	"testing"

	gosync "github.com/ovotech/go-sync"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	slackClient := newMockISlackConversation(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	assert.Equal(t, "test", adapter.conversationName)
	assert.False(t, adapter.MuteRestrictedErrOnKickFromPublic)
	assert.Zero(t, slackClient.Calls)
}

func TestConversation_Get(t *testing.T) {
	t.Parallel()

	slackClient := newMockISlackConversation(t)
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

	slackClient := newMockISlackConversation(t)
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
}

func TestConversation_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackConversation(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient
		adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

		slackClient.EXPECT().KickUserFromConversation("test", "foo").Return(nil)
		slackClient.EXPECT().KickUserFromConversation("test", "bar").Return(nil)

		err := adapter.Remove(ctx, []string{"foo@email", "bar@email"})

		assert.NoError(t, err)
	})

	t.Run("Restricted kick from public conversation", func(t *testing.T) {
		t.Parallel()

		restrictedAction := errors.New("restricted_action") //nolint:goerr113

		slackClient := newMockISlackConversation(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient
		adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

		slackClient.EXPECT().KickUserFromConversation("test", "foo").Maybe().Return(restrictedAction)
		slackClient.EXPECT().KickUserFromConversation("test", "bar").Maybe().Return(restrictedAction)

		adapter.MuteRestrictedErrOnKickFromPublic = false

		err := adapter.Remove(ctx, []string{"foo@email", "bar@email"})

		assert.Error(t, err)
		assert.ErrorIs(t, err, restrictedAction)

		adapter.MuteRestrictedErrOnKickFromPublic = true

		err = adapter.Remove(ctx, []string{"foo@email", "bar@email"})

		assert.NoError(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(map[gosync.ConfigKey]string{
			SlackAPIKey:           "test",
			SlackConversationName: "conversation",
		})

		assert.NoError(t, err)
		assert.IsType(t, &Conversation{}, adapter)
		assert.Equal(t, "conversation", adapter.(*Conversation).conversationName)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[gosync.ConfigKey]string{
				SlackConversationName: "conversation",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, SlackAPIKey)
		})

		t.Run("missing name", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[gosync.ConfigKey]string{
				SlackAPIKey: "test",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, SlackConversationName)
		})
	})
}
