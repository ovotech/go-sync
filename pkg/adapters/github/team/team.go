package team

import (
	"context"
	"fmt"

	"github.com/google/go-github/v47/github"
)

// GitHubDiscovery is required because there are multiple ways to convert a GitHub email into a username.
// At OVO we use SAML, but other organisations may use public emails or another mechanism.
type GitHubDiscovery interface {
	GetUsernameFromEmail(context.Context, []string) ([]string, error)
	GetEmailFromUsername(context.Context, []string) ([]string, error)
}

type iGitHubTeam interface {
	ListTeamMembersBySlug(
		context.Context,
		string,
		string,
		*github.TeamListTeamMembersOptions,
	) ([]*github.User, *github.Response, error)
	AddTeamMembershipBySlug(
		context.Context,
		string,
		string,
		string,
		*github.TeamAddTeamMembershipOptions,
	) (*github.Membership, *github.Response, error)
	RemoveTeamMembershipBySlug(
		context.Context,
		string,
		string,
		string,
	) (*github.Response, error)
}

type Team struct {
	client    iGitHubTeam       // GitHub v3 REST API client.
	discovery GitHubDiscovery   // Discovery adapter to convert GH users -> emails (and vice versa).
	org       string            // GitHub organisation.
	slug      string            // GitHub team slug.
	cache     map[string]string // Cache of users.
}

// New instantiates a new GitHub Team adapter.
func New(client *github.Client, discovery GitHubDiscovery, org string, slug string, optsFn ...func(*Team)) *Team {
	team := &Team{
		client:    client.Teams,
		discovery: discovery,
		org:       org,
		slug:      slug,
		cache:     make(map[string]string),
	}

	for _, fn := range optsFn {
		fn(team)
	}

	return team
}

// Get a list of emails in a GitHub team.
func (t *Team) Get(ctx context.Context) ([]string, error) {
	var out []string

	opts := &github.TeamListTeamMembersOptions{}

	for {
		users, resp, err := t.client.ListTeamMembersBySlug(ctx, t.org, t.slug, opts)
		if err != nil {
			return nil, fmt.Errorf("github.team.get.listteammembersbyslug(%s, %s) -> %w", t.org, t.slug, err)
		}

		logins := make([]string, 0, len(users))
		for _, user := range users {
			logins = append(logins, *user.Login)
		}

		emails, err := t.discovery.GetEmailFromUsername(ctx, logins)
		if err != nil {
			return nil, fmt.Errorf("github.team.get.discovery -> %w", err)
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

// Add emails to a GitHub Team.
func (t *Team) Add(ctx context.Context, emails []string) error {
	names, err := t.discovery.GetUsernameFromEmail(ctx, emails)
	if err != nil {
		return fmt.Errorf("github.team.add.discovery -> %w", err)
	}

	for _, name := range names {
		var opts = &github.TeamAddTeamMembershipOptions{
			Role: "member",
		}

		_, _, err = t.client.AddTeamMembershipBySlug(ctx, t.org, t.slug, name, opts)
		if err != nil {
			return fmt.Errorf("github.team.add.addteammembershipbyslug(%s, %s, %s) -> %w", t.org, t.slug, name, err)
		}
	}

	return nil
}

// Remove a list of emails from a GitHub team.
func (t *Team) Remove(ctx context.Context, emails []string) error {
	if len(t.cache) == 0 {
		if _, err := t.Get(ctx); err != nil {
			return fmt.Errorf("github.team.remove.get -> %w", err)
		}
	}

	for _, email := range emails {
		name := t.cache[email]

		_, err := t.client.RemoveTeamMembershipBySlug(ctx, t.org, t.slug, name)
		if err != nil {
			return fmt.Errorf("github.team.remove.removeteammembershipbyslug -> %w", err)
		}
	}

	return nil
}
