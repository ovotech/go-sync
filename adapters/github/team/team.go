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

# Examples

See [New] and [Init].
*/
package team

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v47/github"
	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/github/discovery"
	"github.com/ovotech/go-sync/adapters/github/discovery/saml"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

/*
GitHubToken is the token used to authenticate with GitHub.
See package docs for more information on how to obtain this token.
*/
const GitHubToken gosync.ConfigKey = "github_token"

/*
GitHubOrg is the name of your GitHub organisation.

https://docs.github.com/en/organizations/collaborating-with-groups-in-organizations/about-organizations

For example:

	https://github.com/ovotech/go-sync

`ovotech` is the name of our organisation.
*/
const GitHubOrg gosync.ConfigKey = "github_org"

/*
TeamSlug is the name of your team slug within your organisation.

For example:

	https://github.com/orgs/ovotech/teams/foobar

`foobar` is the name of our team slug.
*/
const TeamSlug gosync.ConfigKey = "team_slug"

/*
DiscoveryMechanism for converting emails into GitHub users and vice versa. Supported values are:
  - [saml]
*/
const DiscoveryMechanism gosync.ConfigKey = "discovery_mechanism"

var (
	_ gosync.Adapter = &Team{} // Ensure [team.Team] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn  = Init    // Ensure the [team.Init] function fully satisfies the [gosync.InitFn] type.
)

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

type Team struct {
	teams     iGitHubTeam               // GitHub v3 REST API teams.
	discovery discovery.GitHubDiscovery // DiscoveryMechanism adapter to convert GH users -> emails (and vice versa).
	org       string                    // GitHub organisation.
	slug      string                    // GitHub team slug.
	cache     map[string]string         // Cache of users.
	Logger    *log.Logger
}

// Get email addresses in a GitHub Team.
func (t *Team) Get(ctx context.Context) ([]string, error) {
	t.Logger.Printf("Fetching accounts from GitHub team %s/%s", t.org, t.slug)

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

	t.Logger.Println("Fetched accounts successfully")

	return out, nil
}

// Add email addresses to a GitHub Team.
func (t *Team) Add(ctx context.Context, emails []string) error {
	t.Logger.Printf("Adding %s to GitHub team %s/%s", emails, t.org, t.slug)

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

	t.Logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove email addresses from a GitHub Team.
func (t *Team) Remove(ctx context.Context, emails []string) error {
	t.Logger.Printf("Removing %s from GitHub team %s/%s", emails, t.org, t.slug)

	if t.cache == nil {
		return fmt.Errorf("github.team.remove -> %w", gosync.ErrCacheEmpty)
	}

	for _, email := range emails {
		name := t.cache[email]

		_, err := t.teams.RemoveTeamMembershipBySlug(ctx, t.org, t.slug, name)
		if err != nil {
			return fmt.Errorf("github.team.remove.removeteammembershipbyslug -> %w", err)
		}
	}

	t.Logger.Println("Finished removing accounts successfully")

	return nil
}

/*
New GitHub Team [gosync.Adapter].

Recommended reading for parameters:
  - org: [team.GitHubOrg]
  - slug: [team.TeamSlug]
*/
func New(
	client *github.Client,
	discovery discovery.GitHubDiscovery,
	org string,
	slug string,
	optsFn ...func(*Team),
) *Team {
	team := &Team{
		teams:     client.Teams,
		discovery: discovery,
		org:       org,
		slug:      slug,
		cache:     nil,
		Logger:    log.New(os.Stderr, "[go-sync/github/team] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(team)
	}

	return team
}

/*
Init a new GitHub Team [gosync.Adapter].

Required config:
  - [team.GitHubToken]
  - [team.GitHubOrg]
  - [team.TeamSlug]
  - [team.DiscoveryMechanism]
*/
func Init(ctx context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{GitHubToken, GitHubOrg, TeamSlug, DiscoveryMechanism} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("github.team.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config[GitHubToken]},
	))

	var (
		gitHubV3Client = github.NewClient(oauthClient)
		gitHubV4Client = githubv4.NewClient(oauthClient)
		discoverySvc   discovery.GitHubDiscovery
	)

	switch config[DiscoveryMechanism] {
	case "saml":
		discoverySvc = saml.New(gitHubV4Client, config[GitHubOrg])
	default:
		return nil, fmt.Errorf("github.team.init -> %w(%s)", gosync.ErrInvalidConfig, config[DiscoveryMechanism])
	}

	return New(gitHubV3Client, discoverySvc, config[GitHubOrg], config[TeamSlug]), nil
}
