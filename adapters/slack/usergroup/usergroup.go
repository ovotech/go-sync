/*
Package usergroup synchronises emails with a Slack UserGroup.

# Warning

The Slack usergroup API doesn't allow a usergroup to have no members. If this behaviour is expected, we recommend
setting `adapter.MuteGroupCannotBeEmpty = true` to mute the error. No members will be removed, but Go Sync will continue
processing.

# Obtaining the Slack UserGroup ID

In order to use this adapter, you'll need an authenticated Slack client and the ID of the usergroup. This isn't
particularly easy to find, you'll need to log in to Slack via a web browser, and navigate to `People & User Groups`.
Find your UserGroup, and the URL will contain the ID of the group.

For example:

	https://app.slack.com/client/FOOBAR123/browse-user-groups/user_groups/S0123ABC456

`S0123ABC456` is the UserGroup ID.

# Requirements

In order to synchronise with Slack, you'll need to [create a Slack app] with the following OAuth Bot Token permissions:
  - [users:read]
  - [users:read.email]
  - [usergroups:read]
  - [usergroups:write]

# Examples

See [New] and [Init].

[create a Slack app]: https://api.slack.com/authentication/basics
[users:read]: https://api.slack.com/scopes/users:read
[users:read.email]: https://api.slack.com/scopes/users:read.email
[usergroups:read]: https://api.slack.com/scopes/usergroups:read
[usergroups:write]: https://api.slack.com/scopes/usergroups:write
*/
package usergroup

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

// UserGroupID is the Slack UserGroup ID.
const UserGroupID gosync.ConfigKey = "usergroup_id"

// MuteGroupCannotBeEmpty silences errors when removing all users from a UserGroup.
const MuteGroupCannotBeEmpty gosync.ConfigKey = "mute_group_cannot_be_empty"

var (
	_ gosync.Adapter = &UserGroup{} // Ensure [usergroup.UserGroup] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn  = Init         // Ensure the [usergroup.Init] function fully satisfies the [gosync.InitFn] type.
)

// iSlackUserGroup is a subset of the Slack Client, and used to build mocks for easy testing.
type iSlackUserGroup interface {
	GetUserGroupMembersContext(ctx context.Context, userGroup string) ([]string, error)
	GetUsersInfoContext(ctx context.Context, users ...string) (*[]slack.User, error)
	GetUserByEmailContext(ctx context.Context, email string) (*slack.User, error)
	UpdateUserGroupMembersContext(ctx context.Context, userGroup string, members string) (slack.UserGroup, error)
}

type UserGroup struct {
	client      iSlackUserGroup
	userGroupID string
	cache       map[string]string
	Logger      *log.Logger

	MuteGroupCannotBeEmpty bool // See [usergroup.MuteGroupCannotBeEmpty]
}

// Get email addresses in a Slack UserGroup.
func (u *UserGroup) Get(ctx context.Context) ([]string, error) {
	u.Logger.Printf("Fetching accounts from Slack UserGroup %s", u.userGroupID)

	// Initialise the cache.
	u.cache = make(map[string]string)

	// Retrieve a plain list of Slack IDs in the UserGroup.
	groupMembers, err := u.client.GetUserGroupMembersContext(ctx, u.userGroupID)
	if err != nil {
		return nil, fmt.Errorf("slack.usergroup.get.getusergroupmembers -> %w", err)
	}

	// Get the user info for each of the users.
	users, err := u.client.GetUsersInfoContext(ctx, groupMembers...)
	if err != nil {
		return nil, fmt.Errorf("slack.usergroup.get.getusersinfo -> %w", err)
	}

	emails := make([]string, 0, len(*users))
	for _, user := range *users {
		emails = append(emails, user.Profile.Email)
		// Add the email -> ID map for use with the AddAccount and RemoveAccount methods.
		u.cache[user.Profile.Email] = user.ID
	}

	u.Logger.Println("Fetched accounts successfully")

	return emails, nil
}

// Add email addresses to a Slack UserGroup.
func (u *UserGroup) Add(ctx context.Context, emails []string) error {
	u.Logger.Printf("Adding %s to Slack UserGroup %s", emails, u.userGroupID)

	if u.cache == nil {
		return fmt.Errorf("slack.usergroup.add -> %w", gosync.ErrCacheEmpty)
	}

	// The updatedUserGroup is existing users + new users.
	updatedUserGroup := make([]string, 0, len(u.cache)+len(emails))

	// Prefill the updatedUserGroup with everyone currently in the group.
	for _, id := range u.cache {
		updatedUserGroup = append(updatedUserGroup, id)
	}

	// Loop over the emails to be added, and retrieve the Slack IDs.
	for _, email := range emails {
		user, err := u.client.GetUserByEmailContext(ctx, email)
		if err != nil {
			return fmt.Errorf("slack.usergroup.add.getuserbyemail(%s) -> %w", email, err)
		}
		// Add the new email user IDs to the list.
		updatedUserGroup = append(updatedUserGroup, user.ID)

		// Calls to GetUserByEmail are heavily rate limited, so sleep to avoid this.
		time.Sleep(2 * time.Second) //nolint:gomnd
	}

	// Add the members to the Slack UserGroup.
	joinedSlackIds := strings.Join(updatedUserGroup, ",")

	_, err := u.client.UpdateUserGroupMembersContext(ctx, u.userGroupID, joinedSlackIds)
	if err != nil {
		return fmt.Errorf("slack.usergroup.add.updateusergroupmembers(%s) -> %w", u.userGroupID, err)
	}

	u.Logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove email addresses from a Slack UserGroup.
func (u *UserGroup) Remove(ctx context.Context, emails []string) error {
	u.Logger.Printf("Removing %s from Slack UserGroup %s", emails, u.userGroupID)

	if u.cache == nil {
		return fmt.Errorf("slack.usergroup.remove -> %w", gosync.ErrCacheEmpty)
	}

	// Convert the list of email addresses into a map to efficiently lookup emails to remove.
	mapOfEmailsToRemove := map[string]bool{}
	for _, email := range emails {
		mapOfEmailsToRemove[email] = true
	}

	// Iterate over the cached map of emails to Slack IDs, and only include those that aren't in the removal map.
	updatedUserGroup := make([]string, 0, len(u.cache)-len(emails))

	for email, slackID := range u.cache {
		// Only include the Slack ID if it's not in the map of emails to remove.
		if !mapOfEmailsToRemove[email] {
			updatedUserGroup = append(updatedUserGroup, slackID)
		}
	}

	// Update the list of members in the UserGroup.
	concatUserList := strings.Join(updatedUserGroup, ",")

	_, err := u.client.UpdateUserGroupMembersContext(ctx, u.userGroupID, concatUserList)
	if err != nil {
		if strings.Contains(err.Error(), "invalid_arguments") && u.MuteGroupCannotBeEmpty {
			u.Logger.Println("Cannot remove all members from usergroup, but error is muted by configuration - continuing")

			return nil
		}

		return fmt.Errorf("slack.usergroup.remove.updateusergroupmembers(%s, ...) -> %w", u.userGroupID, err)
	}

	u.Logger.Println("Finished removing accounts successfully")

	return nil
}

// New Slack UserGroup [gosync.Adapter].
func New(slackClient *slack.Client, userGroupID string, optsFn ...func(group *UserGroup)) *UserGroup {
	ugAdapter := &UserGroup{
		client:                 slackClient,
		userGroupID:            userGroupID,
		cache:                  nil,
		MuteGroupCannotBeEmpty: false,
		Logger:                 log.New(os.Stderr, "[go-sync/slack/usergroup] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(ugAdapter)
	}

	return ugAdapter
}

/*
Init a new Slack UserGroup [gosync.Adapter].

Required config:
  - [usergroup.SlackAPIKey]
  - [usergroup.UserGroupID]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{SlackAPIKey, UserGroupID} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("slack.conversation.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	client := slack.New(config[SlackAPIKey])

	adapter := New(client, config[UserGroupID])

	if val, ok := config[MuteGroupCannotBeEmpty]; ok {
		adapter.MuteGroupCannotBeEmpty = strings.ToLower(val) == "true"
	}

	return adapter, nil
}
