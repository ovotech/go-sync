package conversation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"

	gosync "github.com/ovotech/go-sync"
)

func TestNew(t *testing.T) {
	t.Parallel()

	slackClient := newMockISlackConversation(t)

	adapter := &Conversation{
		client:           slackClient,
		conversationName: "test",
		Logger:           log.New(os.Stdout, "", log.LstdFlags),
	}

	assert.Equal(t, "test", adapter.conversationName)
	assert.False(t, adapter.MuteRestrictedErrOnKickFromPublic)
	assert.Zero(t, slackClient.Calls)
}

func TestConversation_Get(t *testing.T) {
	t.Parallel()

	slackClient := newMockISlackConversation(t)

	adapter := &Conversation{
		client:           slackClient,
		conversationName: "test",
		Logger:           log.New(os.Stdout, "", log.LstdFlags),
	}

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

func TestConversation_Get_Pagination(t *testing.T) {
	t.Parallel()

	slackClient := newMockISlackConversation(t)

	adapter := &Conversation{
		client:           slackClient,
		conversationName: "test",
		Logger:           log.New(os.Stdout, "", log.LstdFlags),
	}

	incrementingSlice := make([]string, 60)
	firstPage := make([]interface{}, 30)
	secondPage := make([]interface{}, 30)
	firstResponse := make([]slack.User, 30)
	secondResponse := make([]slack.User, 30)

	for idx := range incrementingSlice {
		incrementingSlice[idx] = fmt.Sprint(idx)

		if idx < 30 {
			firstPage[idx] = fmt.Sprint(idx)
			firstResponse[idx] = slack.User{
				ID: fmt.Sprint(idx), IsBot: false, Profile: slack.UserProfile{Email: fmt.Sprint(idx)},
			}
		} else {
			secondPage[idx-30] = fmt.Sprint(idx)
			secondResponse[idx-30] = slack.User{
				ID: fmt.Sprint(idx), IsBot: false, Profile: slack.UserProfile{Email: fmt.Sprint(idx)},
			}
		}
	}

	slackClient.EXPECT().GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: "test",
		Cursor:    "",
		Limit:     50,
	}).Return(incrementingSlice, "", nil)

	slackClient.EXPECT().GetUsersInfo(firstPage...).Return(&firstResponse, nil)
	slackClient.EXPECT().GetUsersInfo(secondPage...).Return(&secondResponse, nil)

	_, err := adapter.Get(context.TODO())

	assert.NoError(t, err)
}

func TestConversation_Add(t *testing.T) {
	t.Parallel()

	slackClient := newMockISlackConversation(t)

	adapter := &Conversation{
		client:           slackClient,
		conversationName: "test",
		Logger:           log.New(os.Stdout, "", log.LstdFlags),
	}

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

//nolint:funlen
func TestConversation_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackConversation(t)

		adapter := &Conversation{
			client:           slackClient,
			conversationName: "test",
			cache:            map[string]string{"foo@email": "foo", "bar@email": "bar"},
			Logger:           log.New(os.Stdout, "", log.LstdFlags),
		}

		slackClient.EXPECT().KickUserFromConversation("test", "foo").Return(nil)
		slackClient.EXPECT().KickUserFromConversation("test", "bar").Return(nil)

		err := adapter.Remove(ctx, []string{"foo@email", "bar@email"})

		assert.NoError(t, err)
	})

	t.Run("Restricted kick from public conversation", func(t *testing.T) {
		t.Parallel()

		restrictedAction := errors.New("restricted_action") //nolint:goerr113

		slackClient := newMockISlackConversation(t)
		adapter := &Conversation{
			client:           slackClient,
			conversationName: "test",
			cache:            map[string]string{"foo@email": "foo", "bar@email": "bar"},
			Logger:           log.New(os.Stdout, "", log.LstdFlags),
		}

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

	t.Run("Check case sensitivity", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackConversation(t)

		adapter := &Conversation{
			client:           slackClient,
			conversationName: "test",
			Logger:           log.New(os.Stdout, "", log.LstdFlags),
		}

		slackClient.EXPECT().GetUsersInConversation(&slack.GetUsersInConversationParameters{
			ChannelID: "test",
			Cursor:    "",
			Limit:     50,
		}).Return([]string{"foo", "bar"}, "", nil)

		slackClient.EXPECT().GetUsersInfo("foo", "bar").Return(&[]slack.User{
			// Capitalise the letter E in email.
			{ID: "foo", IsBot: false, Profile: slack.UserProfile{Email: "foo@Email"}},
			{ID: "bar", IsBot: false, Profile: slack.UserProfile{Email: "bar@Email"}},
		}, nil)

		_, _ = adapter.Get(ctx)

		slackClient.EXPECT().KickUserFromConversation("test", "foo").Return(nil)
		slackClient.EXPECT().KickUserFromConversation("test", "bar").Return(nil)

		// Capitalise the first letter of each email.
		err := adapter.Remove(ctx, []string{"Foo@email", "Bar@email"})

		assert.NoError(t, err)
	})
}

//nolint:funlen
func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			SlackAPIKey: "test",
			Name:        "conversation",
		})

		assert.NoError(t, err)
		assert.IsType(t, &Conversation{}, adapter)
		assert.Equal(t, "conversation", adapter.conversationName)
		assert.False(t, adapter.MuteRestrictedErrOnKickFromPublic)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				Name: "conversation",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, SlackAPIKey)
		})

		t.Run("missing name", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				SlackAPIKey: "test",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, Name)
		})

		t.Run("MuteRestrictedErrOnKickFromPublic", func(t *testing.T) {
			t.Parallel()

			for _, test := range []string{"", "false", "FALSE", "False", "foobar", "test"} {
				adapter, err := Init(ctx, map[gosync.ConfigKey]string{
					SlackAPIKey:                       "test",
					Name:                              "conversation",
					MuteRestrictedErrOnKickFromPublic: test,
				})

				assert.NoError(t, err)
				assert.False(t, adapter.MuteRestrictedErrOnKickFromPublic, test)
			}

			for _, test := range []string{"true", "True", "TRUE"} {
				adapter, err := Init(ctx, map[gosync.ConfigKey]string{
					SlackAPIKey:                       "test",
					Name:                              "conversation",
					MuteRestrictedErrOnKickFromPublic: test,
				})

				assert.NoError(t, err)
				assert.True(t, adapter.MuteRestrictedErrOnKickFromPublic, test)
			}
		})
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(os.Stderr, "custom logger", log.LstdFlags)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			SlackAPIKey: "test",
			Name:        "conversation",
		}, WithLogger(logger))

		assert.NoError(t, err)
		assert.Equal(t, logger, adapter.Logger)
	})

	t.Run("with client", func(t *testing.T) {
		t.Parallel()

		client := slack.New("test")

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Name: "conversation",
		}, WithClient(client))

		assert.NoError(t, err)
		assert.Equal(t, client, adapter.client)
	})
}
