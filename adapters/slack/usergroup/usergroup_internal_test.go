package usergroup

import (
	"context"
	"errors"
	"strings"
	"testing"

	gosync "github.com/ovotech/go-sync"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	slackClient := newMockISlackUserGroup(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	assert.Equal(t, "test", adapter.userGroupName)
	assert.False(t, adapter.MuteGroupCannotBeEmpty)
	assert.Zero(t, slackClient.Calls)
}

func TestUserGroup_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	slackClient := newMockISlackUserGroup(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	slackClient.EXPECT().GetUserGroupMembersContext(ctx, "test").Return([]string{"foo", "bar"}, nil)
	slackClient.EXPECT().GetUsersInfoContext(ctx, "foo", "bar").Maybe().Return(&[]slack.User{
		{ID: "foo", Profile: slack.UserProfile{Email: "foo@email"}},
		{ID: "bar", Profile: slack.UserProfile{Email: "bar@email"}},
	}, nil)
	slackClient.EXPECT().GetUsersInfoContext(ctx, "bar", "foo").Maybe().Return(&[]slack.User{
		{ID: "bar", Profile: slack.UserProfile{Email: "bar@email"}},
		{ID: "foo", Profile: slack.UserProfile{Email: "foo@email"}},
	}, nil)

	users, err := adapter.Get(ctx)

	assert.NoError(t, err)
	assert.Equal(t, []string{"foo@email", "bar@email"}, users)
}

func TestUserGroup_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("No cache", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient

		err := adapter.Add(ctx, []string{"foo", "bar"})

		assert.Error(t, err)
		assert.ErrorIs(t, err, gosync.ErrCacheEmpty)
	})

	t.Run("Add accounts", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient

		slackClient.EXPECT().GetUserByEmailContext(ctx, "fizz@email").Return(&slack.User{ID: "fizz"}, nil)
		slackClient.EXPECT().GetUserByEmailContext(ctx, "buzz@email").Return(&slack.User{ID: "buzz"}, nil)
		slackClient.EXPECT().UpdateUserGroupMembersContext(ctx,
			"test", mock.Anything,
		).Run(func(_ context.Context, userGroup string, members string) { //nolint:contextcheck
			assert.Equal(t, "test", userGroup)
			assert.ElementsMatch(t, strings.Split(members, ","), []string{"foo", "bar", "fizz", "buzz"})
		}).Return(slack.UserGroup{DateDelete: 0}, nil)

		adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}
		err := adapter.Add(ctx, []string{"fizz@email", "buzz@email"})

		assert.NoError(t, err)
	})
}

func TestUserGroup_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("No cache", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient

		err := adapter.Remove(ctx, []string{"foo@email"})

		assert.Error(t, err)
		assert.ErrorIs(t, err, gosync.ErrCacheEmpty)
	})

	t.Run("Remove accounts", func(t *testing.T) {
		t.Parallel()

		slackClient := newMockISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient
		adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

		slackClient.EXPECT().UpdateUserGroupMembersContext(ctx, "test", "foo").Return(slack.UserGroup{}, nil)

		err := adapter.Remove(ctx, []string{"bar@email"})

		assert.NoError(t, err)
	})

	t.Run("Return/mute error if number of accounts reaches zero", func(t *testing.T) {
		t.Parallel()

		// Mock the error returned from the Slack API.
		errInvalidArguments := errors.New("invalid_arguments") //nolint:goerr113

		slackClient := newMockISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient
		adapter.cache = map[string]string{"foo@email": "foo"}
		adapter.MuteGroupCannotBeEmpty = false

		slackClient.EXPECT().UpdateUserGroupMembersContext(ctx, "test", "").Return(slack.UserGroup{}, errInvalidArguments)

		err := adapter.Remove(ctx, []string{"foo@email"})

		assert.ErrorIs(t, err, errInvalidArguments)

		// Reset the cache and mute the empty group error.
		adapter.MuteGroupCannotBeEmpty = true

		err = adapter.Remove(ctx, []string{"foo@email"})

		assert.NoError(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(map[gosync.ConfigKey]string{
			SlackAPIKey:        "test",
			SlackUserGroupName: "usergroup",
		})

		assert.NoError(t, err)
		assert.IsType(t, &UserGroup{}, adapter)
		assert.Equal(t, "usergroup", adapter.(*UserGroup).userGroupName)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[gosync.ConfigKey]string{
				SlackUserGroupName: "usergroup",
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
			assert.ErrorContains(t, err, SlackUserGroupName)
		})
	})
}
