/*
Package team synchronises email addresses with GitHub teams.

# Discovery

You must provide a [discovery] adapter in order to use this adapter. This is because converting email addresses to
GitHub usernames isn't straightforward. At OVO, we enforce SAML for our GitHub users, and have provided a
SAML -> GitHub Username discovery adapter, but you may need to write your own.

# Requirements

In order to synchronise with GitHub, you'll need to create a [Personal Access Token] with the following permissions:
  - admin:org
  - write:org
  - read:org

https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
*/
package team

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v47/github"

	"github.com/ovotech/go-sync/packages/github/discovery"
	"github.com/ovotech/go-sync/packages/gosync"
)

var _ gosync.Adapter = &Adapter{} // Ensure [team.Adapter] fully satisfies the [gosync.Adapter] interface.

// iSlackConversation is a subset of the Slack Client used to build mocks for easy testing.
type iGitHubTeam interface {
	ListTeamMembersBySlug(
		ctx context.Context,
		org,
		slug string,
		opts *github.TeamListTeamMembersOptions,
	) ([]*github.User, *github.Response, error)
	AddTeamMembershipBySlug(
		ctx context.Context,
		org,
		slug,
		user string,
		opts *github.TeamAddTeamMembershipOptions,
	) (*github.Membership, *github.Response, error)
	RemoveTeamMembershipBySlug(ctx context.Context, org, slug, user string) (*github.Response, error)
}

type Adapter struct {
	teams     iGitHubTeam               // GitHub v3 REST API teams.
	discovery discovery.GitHubDiscovery // DiscoveryMechanism adapter to convert GH users -> emails (and vice versa).
	org       string                    // GitHub organisation.
	slug      string                    // GitHub team slug.
	cache     map[string]string         // Cache of users.
	logger    *log.Logger
}

// Get email addresses in a GitHub Adapter.
func (a *Adapter) Get(ctx context.Context) ([]string, error) {
	a.logger.Printf("Fetching accounts from GitHub team %s/%s", a.org, a.slug)

	// Initialise the cache.
	a.cache = make(map[string]string)

	out := make([]string, 0)

	opts := &github.TeamListTeamMembersOptions{}

	for {
		users, resp, err := a.teams.ListTeamMembersBySlug(ctx, a.org, a.slug, opts)
		if err != nil {
			return nil, fmt.Errorf("github.team.get.listteammembersbyslug(%s, %s) -> %w", a.org, a.slug, err)
		}

		logins := make([]string, 0, len(users))
		for _, user := range users {
			logins = append(logins, *user.Login)
		}

		emails, err := a.discovery.GetEmailFromUsername(ctx, logins)
		if err != nil {
			return nil, fmt.Errorf("github.team.get.discovery -> %w", err)
		}

		out = append(out, emails...)

		for index, user := range users {
			a.cache[emails[index]] = *user.Login
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	a.logger.Println("Fetched accounts successfully")

	return out, nil
}

// Add email addresses to a GitHub Adapter.
func (a *Adapter) Add(ctx context.Context, emails []string) error {
	a.logger.Printf("Adding %s to GitHub team %s/%s", emails, a.org, a.slug)

	names, err := a.discovery.GetUsernameFromEmail(ctx, emails)
	if err != nil {
		return fmt.Errorf("github.team.add.discovery -> %w", err)
	}

	for _, name := range names {
		opts := &github.TeamAddTeamMembershipOptions{
			Role: "member",
		}

		_, _, err = a.teams.AddTeamMembershipBySlug(ctx, a.org, a.slug, name, opts)
		if err != nil {
			return fmt.Errorf("github.team.add.addteammembershipbyslug(%s, %s, %s) -> %w", a.org, a.slug, name, err)
		}
	}

	a.logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove email addresses from a GitHub Adapter.
func (a *Adapter) Remove(ctx context.Context, emails []string) error {
	a.logger.Printf("Removing %s from GitHub team %s/%s", emails, a.org, a.slug)

	if a.cache == nil {
		return fmt.Errorf("github.team.remove -> %w", gosync.ErrCacheEmpty)
	}

	for _, email := range emails {
		name := a.cache[email]

		_, err := a.teams.RemoveTeamMembershipBySlug(ctx, a.org, a.slug, name)
		if err != nil {
			return fmt.Errorf("github.team.remove.removeteammembershipbyslug -> %w", err)
		}
	}

	a.logger.Println("Finished removing accounts successfully")

	return nil
}
