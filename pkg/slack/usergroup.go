package slack

import (
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

type UserGroup struct {
	client        *slack.Client
	userGroupName string
	cache         map[string]string
}

// NewUserGroupService instantiates a new Slack UserGroup service.
func NewUserGroupService(slackClient *slack.Client, userGroup string) *UserGroup {
	return &UserGroup{
		client:        slackClient,
		userGroupName: userGroup,
		cache:         map[string]string{},
	}
}

func (u *UserGroup) get() ([]string, error) {
	// Retrieve a plain list of Slack IDs in the user group.
	groupMembers, err := u.client.GetUserGroupMembers(u.userGroupName)
	if err != nil {
		return nil, fmt.Errorf("GetUserGroupMembers -> %w", err)
	}

	// Get the user info for each of the users.
	users, err := u.client.GetUsersInfo(groupMembers...)
	if err != nil {
		return nil, fmt.Errorf("GetUsersInfo -> %w", err)
	}

	emails := make([]string, 0, len(*users))
	for _, user := range *users {
		emails = append(emails, user.Profile.Email)
		// Add the email -> ID map for use with the AddAccount and RemoveAccount methods.
		u.cache[user.Profile.Email] = user.ID
	}

	return emails, nil
}

// Get a list of email addresses in a Slack User Group.
func (u *UserGroup) Get() ([]string, error) {
	ids, err := u.get()
	if err != nil {
		return nil, fmt.Errorf("slack.usergroup.Get -> %w", err)
	}

	return ids, nil
}

// Add a list of emails to a Slack UserGroup.
// Since the Slack API takes all of this as a single request, it either returns a full list of successful emails,
// or an error.
func (u *UserGroup) Add(emails ...string) ([]string, []error, error) {
	// If the cache hasn't been generated, regenerate it.
	if len(u.cache) == 0 {
		_, err := u.get()
		if err != nil {
			return nil, nil, fmt.Errorf("slack.usergroup.Add.get -> %w", err)
		}
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
			return nil, nil, fmt.Errorf("slack.usergroup.Add.GetUserByEmail(%s) -> %w", email, err)
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

		return nil, nil, fmt.Errorf("slack.usergroup.Add.UpdateUserGroupMembers(%s) -> %w", u.userGroupName, err)
	}

	// If the group is disabled, re-enable it.
	if group.DateDelete != 0 {
		_, err = u.client.EnableUserGroup(u.userGroupName)
		if err != nil {
			return nil, nil, fmt.Errorf("slack.usergroup.Add.EnableUserGroup(%s) -> %w", u.userGroupName, err)
		}
	}

	return emails, nil, nil
}

// Remove a list of email addresses from a Slack UserGroup.
func (u *UserGroup) Remove(emails ...string) ([]string, []error, error) {
	// If the cache hasn't been generated, regenerate it.
	if len(u.cache) == 0 {
		_, err := u.get()
		if err != nil {
			return nil, nil, fmt.Errorf("slack.usergroup.Remove.get -> %w", err)
		}
	}

	// If this change would remove all users, disable the group instead.
	if len(u.cache) == len(emails) {
		_, err := u.client.DisableUserGroup(u.userGroupName)
		if err != nil {
			return nil, nil, fmt.Errorf("slack.usergroup.Remove.DisableUserGroup(%s) -> %w", u.userGroupName, err)
		}

		return emails, nil, nil
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
		return nil, nil, fmt.Errorf("slack.usergroup.Remove.UpdateUserGroupMembers(%s, ...) -> %w", u.userGroupName, err)
	}

	return emails, nil, nil
}
