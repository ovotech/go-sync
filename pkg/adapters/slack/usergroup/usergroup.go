/*
Package usergroup synchronises emails with a Slack UserGroup.

In order to use this adapter, you'll need an authenticated Slack client and the ID of the usergroup.
This isn't particularly easy to find, you'll need to log in to Slack via a web browser, and navigate to
`People & User Groups`. Find your User group, and the URL will contain the ID of the group.
*/
package usergroup

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ovotech/go-sync/internal/types"
	"github.com/ovotech/go-sync/pkg/ports"
	"github.com/slack-go/slack"
)

// Ensure the adapter type fully satisfies the ports.Adapter interface.
var _ ports.Adapter = &UserGroup{}

// iSlackUserGroup is a subset of the Slack Client, and used to build mocks for easy testing.
type iSlackUserGroup interface {
	GetUserGroupMembersContext(ctx context.Context, userGroup string) ([]string, error)
	GetUsersInfoContext(ctx context.Context, users ...string) (*[]slack.User, error)
	GetUserByEmailContext(ctx context.Context, email string) (*slack.User, error)
	UpdateUserGroupMembersContext(ctx context.Context, userGroup string, members string) (slack.UserGroup, error)
}

type UserGroup struct {
	client        iSlackUserGroup
	userGroupName string
	cache         map[string]string
	logger        types.Logger

	// MuteGroupCannotBeEmpty silences errors when removing everyone from a usergroup.
	MuteGroupCannotBeEmpty bool
}

// ErrCacheEmpty shouldn't realistically be raised unless the adapter is being used outside of Go Sync.
var ErrCacheEmpty = errors.New("cache is empty - run Get()")

func copyCache(originalMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range originalMap {
		newMap[k] = v
	}

	return newMap
}

// WithLogger sets a custom logger.
func WithLogger(logger types.Logger) func(*UserGroup) {
	return func(userGroup *UserGroup) {
		userGroup.logger = logger
	}
}

// New instantiates a new Slack UserGroup adapter.
func New(slackClient *slack.Client, userGroup string, optsFn ...func(group *UserGroup)) *UserGroup {
	ugAdapter := &UserGroup{
		client:                 slackClient,
		userGroupName:          userGroup,
		cache:                  nil,
		MuteGroupCannotBeEmpty: false,
		logger:                 log.New(os.Stderr, "[go-sync/slack/usergroup] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(ugAdapter)
	}

	return ugAdapter
}

// Get emails of Slack users in a User group.
func (u *UserGroup) Get(ctx context.Context) ([]string, error) {
	u.logger.Printf("Fetching accounts from Slack UserGroup %s", u.userGroupName)

	// Initialise the cache.
	u.cache = make(map[string]string)

	// Retrieve a plain list of Slack IDs in the user group.
	groupMembers, err := u.client.GetUserGroupMembersContext(ctx, u.userGroupName)
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

	u.logger.Println("Fetched accounts successfully")

	return emails, nil
}

// Add emails to a Slack User group.
func (u *UserGroup) Add(ctx context.Context, emails []string) error {
	u.logger.Printf("Adding %s to Slack UserGroup %s", emails, u.userGroupName)

	if u.cache == nil {
		return fmt.Errorf("slack.usergroup.add -> %w", ErrCacheEmpty)
	}

	// The updatedUserGroup is existing users + new users.
	updatedUserGroup := make([]string, 0, len(u.cache)+len(emails))

	// Prefill the updatedUserGroup with everyone currently in the group.
	for _, id := range u.cache {
		updatedUserGroup = append(updatedUserGroup, id)
	}

	// Shallow copy the cache.
	localCache := copyCache(u.cache)

	// Loop over the emails to be added, and retrieve the Slack IDs.
	for _, email := range emails {
		user, err := u.client.GetUserByEmailContext(ctx, email)
		if err != nil {
			return fmt.Errorf("slack.usergroup.add.getuserbyemail(%s) -> %w", email, err)
		}
		// Add the new email user IDs to the list.
		updatedUserGroup = append(updatedUserGroup, user.ID)

		// Add the entry to the cache.
		localCache[email] = user.ID

		// Calls to GetUserByEmail are heavily rate limited, so sleep to avoid this.
		time.Sleep(2 * time.Second) //nolint:gomnd
	}

	// Add the members to the Slack user group.
	joinedSlackIds := strings.Join(updatedUserGroup, ",")

	_, err := u.client.UpdateUserGroupMembersContext(ctx, u.userGroupName, joinedSlackIds)
	if err != nil {
		return fmt.Errorf("slack.usergroup.add.updateusergroupmembers(%s) -> %w", u.userGroupName, err)
	}

	// Update the cache after successful additions to the usergroup.
	u.cache = localCache

	u.logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove emails from a Slack User group.
func (u *UserGroup) Remove(ctx context.Context, emails []string) error {
	u.logger.Printf("Removing %s from Slack UserGroup %s", emails, u.userGroupName)

	if u.cache == nil {
		return fmt.Errorf("slack.usergroup.remove -> %w", ErrCacheEmpty)
	}

	// Shallow copy the cache.
	localCache := copyCache(u.cache)

	// Convert the list of email addresses into a map to efficiently lookup emails to remove.
	mapOfEmailsToRemove := map[string]bool{}
	for _, email := range emails {
		mapOfEmailsToRemove[email] = true
	}

	// Iterate over the cached map of emails to Slack IDs, and only include those that aren't in the removal map.
	updatedUserGroup := make([]string, 0, len(localCache)-len(emails))

	for email, slackID := range localCache {
		// Only include the Slack ID if it's not in the map of emails to remove.
		if !mapOfEmailsToRemove[email] {
			updatedUserGroup = append(updatedUserGroup, slackID)
		} else {
			delete(localCache, email)
		}
	}

	// Update the list of members in the user group.
	concatUserList := strings.Join(updatedUserGroup, ",")

	_, err := u.client.UpdateUserGroupMembersContext(ctx, u.userGroupName, concatUserList)
	if err != nil {
		if strings.Contains(err.Error(), "invalid_arguments") && u.MuteGroupCannotBeEmpty {
			u.logger.Println("Cannot remove all members from usergroup, but error is muted by configuration - continuing")

			return nil
		}

		return fmt.Errorf("slack.usergroup.remove.updateusergroupmembers(%s, ...) -> %w", u.userGroupName, err)
	}

	// Update the cache after successful removals from the usergroup.
	u.cache = localCache

	u.logger.Println("Finished removing accounts successfully")

	return nil
}
