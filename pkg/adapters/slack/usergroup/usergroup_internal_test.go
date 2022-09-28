package usergroup

import (
	"context"
	"strings"
	"testing"

	"github.com/ovotech/go-sync/mocks"
	"github.com/ovotech/go-sync/pkg/ports"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImplementsInterface(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*ports.Adapter)(nil), &UserGroup{})
}

func TestNew(t *testing.T) {
	t.Parallel()

	slackClient := mocks.NewISlackUserGroup(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	assert.Equal(t, "test", adapter.userGroupName)
	assert.Zero(t, slackClient.Calls)
}

func TestUserGroup_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	slackClient := mocks.NewISlackUserGroup(t)
	adapter := New(&slack.Client{}, "test")
	adapter.client = slackClient

	slackClient.EXPECT().GetUserGroupMembers("test").Return([]string{"foo", "bar"}, nil)
	slackClient.EXPECT().GetUsersInfo("foo", "bar").Maybe().Return(&[]slack.User{
		{ID: "foo", Profile: slack.UserProfile{Email: "foo@email"}},
		{ID: "bar", Profile: slack.UserProfile{Email: "bar@email"}},
	}, nil)
	slackClient.EXPECT().GetUsersInfo("bar", "foo").Maybe().Return(&[]slack.User{
		{ID: "bar", Profile: slack.UserProfile{Email: "bar@email"}},
		{ID: "foo", Profile: slack.UserProfile{Email: "foo@email"}},
	}, nil)

	users, err := adapter.Get(ctx)

	assert.NoError(t, err)
	assert.Equal(t, []string{"foo@email", "bar@email"}, users)
}

//nolint:funlen
func TestUserGroup_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("No cache", func(t *testing.T) {
		t.Parallel()

		slackClient := mocks.NewISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient

		err := adapter.Add(ctx, []string{"foo", "bar"})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCacheEmpty)
	})

	t.Run("Remove accounts", func(t *testing.T) {
		t.Parallel()

		slackClient := mocks.NewISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient

		argsMatch := func(userGroup string, members string) {
			assert.Equal(t, "test", userGroup)
			assert.ElementsMatch(t, strings.Split(members, ","), []string{"foo", "bar", "fizz", "buzz"})
		}

		slackClient.EXPECT().GetUserByEmail("fizz@email").Return(&slack.User{ID: "fizz"}, nil)
		slackClient.EXPECT().GetUserByEmail("buzz@email").Return(&slack.User{ID: "buzz"}, nil)
		slackClient.EXPECT().UpdateUserGroupMembers(
			"test", mock.Anything,
		).Run(argsMatch).Return(slack.UserGroup{DateDelete: 0}, nil)

		adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}
		err := adapter.Add(ctx, []string{"fizz@email", "buzz@email"})

		assert.NoError(t, err)
		assert.Equal(t, adapter.cache, map[string]string{
			"foo@email": "foo", "bar@email": "bar", "fizz@email": "fizz", "buzz@email": "buzz",
		})
	})

	t.Run("Enable user group if disabled", func(t *testing.T) {
		t.Parallel()

		slackClient := mocks.NewISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient
		adapter.cache = map[string]string{"foo@email": "foo"}

		slackClient.EXPECT().GetUserByEmail("bar@email").Return(&slack.User{ID: "bar"}, nil)
		slackClient.EXPECT().UpdateUserGroupMembers("test", "foo,bar").Maybe().
			Return(slack.UserGroup{DateDelete: 1}, nil)
		slackClient.EXPECT().UpdateUserGroupMembers("test", "bar,foo").Maybe().
			Return(slack.UserGroup{DateDelete: 1}, nil)
		slackClient.EXPECT().EnableUserGroup("test").Return(slack.UserGroup{}, nil)

		err := adapter.Add(ctx, []string{"bar@email"})

		assert.NoError(t, err)
	})
}

func TestUserGroup_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("No cache", func(t *testing.T) {
		t.Parallel()

		slackClient := mocks.NewISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient

		err := adapter.Remove(ctx, []string{"foo@email"})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCacheEmpty)
	})

	t.Run("Remove accounts", func(t *testing.T) {
		t.Parallel()

		slackClient := mocks.NewISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient
		adapter.cache = map[string]string{"foo@email": "foo", "bar@email": "bar"}

		slackClient.EXPECT().UpdateUserGroupMembers("test", "foo").Return(slack.UserGroup{}, nil)

		err := adapter.Remove(ctx, []string{"bar@email"})

		assert.NoError(t, err)
	})

	t.Run("Disable Usergroup if membership would be zero", func(t *testing.T) {
		t.Parallel()

		slackClient := mocks.NewISlackUserGroup(t)
		adapter := New(&slack.Client{}, "test")
		adapter.client = slackClient
		adapter.cache = map[string]string{"foo@email": "foo"}

		slackClient.EXPECT().DisableUserGroup("test").Return(slack.UserGroup{}, nil)

		err := adapter.Remove(ctx, []string{"foo@email"})

		assert.NoError(t, err)
	})
}
