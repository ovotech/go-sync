// Package usergroup synchronises emails with a Slack UserGroup.
//
// In order to use this adapter, you'll need an authenticated Slack client and the ID of the usergroup.
// This isn't particularly easy to find, you'll need to log in to Slack via a web browser, and navigate to
// `People & User Groups`. Find your User group, and the URL will contain the ID of the group.
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
	"github.com/slack-go/slack"
)

// iSlackUserGroup is a subset of the Slack Client, and used to build mocks for easy testing.
type iSlackUserGroup interface {
	GetUserGroupMembers(userGroup string) ([]string, error)
	GetUsersInfo(users ...string) (*[]slack.User, error)
	GetUserByEmail(email string) (*slack.User, error)
	UpdateUserGroupMembers(userGroup string, members string) (slack.UserGroup, error)
	EnableUserGroup(userGroup string) (slack.UserGroup, error)
	DisableUserGroup(userGroup string) (slack.UserGroup, error)
}

type UserGroup struct {
	client        iSlackUserGroup
	userGroupName string
	cache         map[string]string
	logger        types.Logger
}

// ErrCacheEmpty shouldn't realistically be raised unless the adapter is being used outside of Go Sync.
var ErrCacheEmpty = errors.New("cache is empty - run Get()")

// OptionLogger can be used to set a custom logger.
func OptionLogger(logger types.Logger) func(*UserGroup) {
	return func(userGroup *UserGroup) {
		userGroup.logger = logger
	}
}

// New instantiates a new Slack UserGroup adapter.
func New(slackClient *slack.Client, userGroup string, optsFn ...func(group *UserGroup)) *UserGroup {
	ugAdapter := &UserGroup{
		client:        slackClient,
		userGroupName: userGroup,
		cache:         map[string]string{},
		logger:        log.New(os.Stderr, "[go-sync/slack/usergroup] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(ugAdapter)
	}

	return ugAdapter
}

// Get emails of Slack users in a User group.
func (u *UserGroup) Get(_ context.Context) ([]string, error) {
	u.logger.Printf("Fetching accounts from Slack UserGroup %s", u.userGroupName)

	// Retrieve a plain list of Slack IDs in the user group.
	groupMembers, err := u.client.GetUserGroupMembers(u.userGroupName)
	if err != nil {
		return nil, fmt.Errorf("slack.usergroup.get.getusergroupmembers -> %w", err)
	}

	// Get the user info for each of the users.
	users, err := u.client.GetUsersInfo(groupMembers...)
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
func (u *UserGroup) Add(_ context.Context, emails []string) error {
	u.logger.Printf("Adding %s to Slack UserGroup %s", emails, u.userGroupName)

	if len(u.cache) == 0 {
		return fmt.Errorf("slack.usergroup.add -> %w", ErrCacheEmpty)
	}

	// The updatedUserGroup is existing users + new users.
	updatedUserGroup := make([]string, 0, len(u.cache)+len(emails))

	// Prefill the updatedUserGroup with everyone currently in the group.
	for _, id := range u.cache {
		updatedUserGroup = append(updatedUserGroup, id)
	}

	// Loop over the emails to be added, and retrieve the Slack IDs.
	for _, email := range emails {
		user, err := u.client.GetUserByEmail(email)
		if err != nil {
			return fmt.Errorf("slack.usergroup.add.getuserbyemail(%s) -> %w", email, err)
		}
		// Add the new email user IDs to the list.
		updatedUserGroup = append(updatedUserGroup, user.ID)

		// Add the entry to the cache.
		u.cache[email] = user.ID

		// Calls to GetUserByEmail are heavily rate limited, so sleep to avoid this.
		time.Sleep(2 * time.Second) //nolint:gomnd
	}

	// Add the members to the Slack user group.
	joinedSlackIds := strings.Join(updatedUserGroup, ",")

	group, err := u.client.UpdateUserGroupMembers(u.userGroupName, joinedSlackIds)
	if err != nil {
		u.cache = map[string]string{}

		return fmt.Errorf("slack.usergroup.add.updateusergroupmembers(%s) -> %w", u.userGroupName, err)
	}

	// If the group is disabled, re-enable it.
	if group.DateDelete != 0 {
		_, err = u.client.EnableUserGroup(u.userGroupName)
		if err != nil {
			return fmt.Errorf("slack.usergroup.add.enableusergroup(%s) -> %w", u.userGroupName, err)
		}
	}

	u.logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove emails from a Slack User group.
func (u *UserGroup) Remove(_ context.Context, emails []string) error {
	u.logger.Printf("Removing %s from Slack UserGroup %s", emails, u.userGroupName)

	if len(u.cache) == 0 {
		return fmt.Errorf("slack.usergroup.remove -> %w", ErrCacheEmpty)
	}

	// If this change would remove all users, disable the group instead.
	if len(u.cache) == len(emails) {
		_, err := u.client.DisableUserGroup(u.userGroupName)
		if err != nil {
			return fmt.Errorf("slack.usergroup.remove.disableusergroup(%s) -> %w", u.userGroupName, err)
		}

		return nil
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
		} else {
			delete(u.cache, email)
		}
	}

	// Update the list of members in the user group.
	concatUserList := strings.Join(updatedUserGroup, ",")

	_, err := u.client.UpdateUserGroupMembers(u.userGroupName, concatUserList)
	if err != nil {
		return fmt.Errorf("slack.usergroup.remove.updateusergroupmembers(%s, ...) -> %w", u.userGroupName, err)
	}

	u.logger.Println("Finished removing accounts successfully")

	return nil
}
