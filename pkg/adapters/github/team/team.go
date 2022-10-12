/*
Package team synchronises emails with GitHub teams.

You must provide a discovery service in order to use this adapter. This is because converting email addresses to
GitHub usernames isn't straightforward. At OVO, we enforce SAML for our GitHub users, and have provided a
SAML -> GitHub Username discovery service, but you may need to write your own.
*/
package team

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v47/github"
	"github.com/ovotech/go-sync/internal/types"
	gosyncerrors "github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/ports"
)

// Ensure the adapter type fully satisfies the ports.Adapter interface.
var _ ports.Adapter = &Team{}

// GitHubDiscovery is required because there are multiple ways to convert a GitHub email into a username.
// At OVO we use SAML, but other organisations may use public emails or another mechanism.
type GitHubDiscovery interface {
	GetUsernameFromEmail(context.Context, []string) ([]string, error)
	GetEmailFromUsername(context.Context, []string) ([]string, error)
}

// iSlackConversation is a subset of the Slack Client, and used to build mocks for easy testing.
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

type Team struct {
	teams     iGitHubTeam       // GitHub v3 REST API teams.
	discovery GitHubDiscovery   // Discovery adapter to convert GH users -> emails (and vice versa).
	org       string            // GitHub organisation.
	slug      string            // GitHub team slug.
	cache     map[string]string // Cache of users.
	logger    types.Logger
}

// WithLogger sets a custom logger.
func WithLogger(logger types.Logger) func(*Team) {
	return func(team *Team) {
		team.logger = logger
	}
}

// New instantiates a new GitHub Team adapter.
func New(client *github.Client, discovery GitHubDiscovery, org string, slug string, optsFn ...func(*Team)) *Team {
	team := &Team{
		teams:     client.Teams,
		discovery: discovery,
		org:       org,
		slug:      slug,
		cache:     nil,
		logger:    log.New(os.Stderr, "[go-sync/github/team] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(team)
	}

	return team
}

// Get emails of users in a GitHub team.
func (t *Team) Get(ctx context.Context) ([]string, error) {
	t.logger.Printf("Fetching accounts from GitHub team %s/%s", t.org, t.slug)

	// Initialise the cache.
	t.cache = make(map[string]string)

	out := make([]string, 0)

	opts := &github.TeamListTeamMembersOptions{}

	for {
		users, resp, err := t.teams.ListTeamMembersBySlug(ctx, t.org, t.slug, opts)
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

	t.logger.Println("Fetched accounts successfully")

	return out, nil
}

// Add emails to a GitHub Team.
func (t *Team) Add(ctx context.Context, emails []string) error {
	t.logger.Printf("Adding %s to GitHub team %s/%s", emails, t.org, t.slug)

	names, err := t.discovery.GetUsernameFromEmail(ctx, emails)
	if err != nil {
		return fmt.Errorf("github.team.add.discovery -> %w", err)
	}

	for _, name := range names {
		var opts = &github.TeamAddTeamMembershipOptions{
			Role: "member",
		}

		_, _, err = t.teams.AddTeamMembershipBySlug(ctx, t.org, t.slug, name, opts)
		if err != nil {
			return fmt.Errorf("github.team.add.addteammembershipbyslug(%s, %s, %s) -> %w", t.org, t.slug, name, err)
		}
	}

	t.logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove emails from a GitHub team.
func (t *Team) Remove(ctx context.Context, emails []string) error {
	t.logger.Printf("Removing %s from GitHub team %s/%s", emails, t.org, t.slug)

	if t.cache == nil {
		return fmt.Errorf("github.team.remove -> %w", gosyncerrors.ErrCacheEmpty)
	}

	for _, email := range emails {
		name := t.cache[email]

		_, err := t.teams.RemoveTeamMembershipBySlug(ctx, t.org, t.slug, name)
		if err != nil {
			return fmt.Errorf("github.team.remove.removeteammembershipbyslug -> %w", err)
		}
	}

	t.logger.Println("Finished removing accounts successfully")

	return nil
}
