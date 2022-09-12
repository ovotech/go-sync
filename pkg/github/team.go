package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v47/github"
	"github.com/ovotech/go-sync/pkg/core/ports"
)

type Team struct {
	client    *github.Client    // GitHub v3 REST API client.
	discovery ports.Discovery   // Discovery service to convert emails into GH users.
	org       string            // GitHub organisation.
	slug      string            // GitHub team slug.
	cache     map[string]string // Cache of users.
}

// NewTeamSyncService instantiates a new GitHub Team service.
func NewTeamSyncService(client *github.Client, discovery ports.Discovery, org string, slug string) *Team {
	return &Team{
		client:    client,
		discovery: discovery,
		org:       org,
		slug:      slug,
		cache:     map[string]string{},
	}
}

func (t *Team) Get() ([]string, error) {
	var out []string

	opts := &github.TeamListTeamMembersOptions{} //nolint:exhaustruct

	for {
		users, resp, err := t.client.Teams.ListTeamMembersBySlug(context.Background(), t.org, t.slug, opts)
		if err != nil {
			return nil, fmt.Errorf("github.team.Get.ListTeamMembersBySlug(%s, %s) -> %w", t.org, t.slug, err)
		}

		logins := make([]string, 0, len(users))
		for _, user := range users {
			logins = append(logins, *user.Login)
		}

		emails, err := t.discovery.GetEmailFromUsername(logins...)
		if err != nil {
			return nil, fmt.Errorf("github.team.Get.discovery -> %w", err)
		}

		out = append(out, emails...)

		for index, user := range users {
			t.cache[emails[index]] = *user.Login
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}

func (t *Team) Add(emails ...string) ([]string, []error, error) {
	var (
		success []string
		failure []error
	)

	names, err := t.discovery.GetUsernameFromEmail(emails...)
	if err != nil {
		return nil, nil, fmt.Errorf("github.team.Add.discovery -> %w", err)
	}

	for index, name := range names {
		var opts = &github.TeamAddTeamMembershipOptions{
			Role: "member",
		}

		_, _, err := t.client.Teams.AddTeamMembershipBySlug(context.Background(), t.org, t.slug, name, opts)
		if err == nil {
			success = append(success, emails[index])
		} else {
			failure = append(
				failure,
				fmt.Errorf("github.team.Add.AddTeamMembershipBySlug(%s, %s, %s) -> %w", t.org, t.slug, name, err),
			)
		}
	}

	return success, failure, nil
}

func (t *Team) Remove(emails ...string) ([]string, []error, error) {
	var (
		success []string
		failure []error
	)

	if len(t.cache) == 0 {
		if _, err := t.Get(); err != nil {
			return nil, nil, fmt.Errorf("github.team.Remove.get -> %w", err)
		}
	}

	for _, email := range emails {
		name := t.cache[email]

		_, err := t.client.Teams.RemoveTeamMembershipBySlug(context.Background(), t.org, t.slug, name)
		if err == nil {
			success = append(success, email)
		} else {
			failure = append(
				failure,
				fmt.Errorf("github.team.Remove.RemoveTeamMembershipBySlug(%s, %s) -> %w", t.org, t.slug, err),
			)
		}
	}

	return success, failure, nil
}
