package team

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v47/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/ovotech/go-sync/packages/github/discovery"
	"github.com/ovotech/go-sync/packages/github/discovery/saml"
	"github.com/ovotech/go-sync/packages/gosync"
)

var _ gosync.InitFn = Init // Ensure the [team.Init] function fully satisfies the [gosync.InitFn] type.

// WithGitHubV3Client passes a custom GitHub V3 client for authentication.
func WithGitHubV3Client(client *github.Client) func(interface{}) {
	return func(i interface{}) {
		i.(*Adapter).teams = client.Teams
	}
}

// WithDiscoveryService sets a custom GitHub [discovery.GitHubDiscovery] service.
func WithDiscoveryService(svc discovery.GitHubDiscovery) func(interface{}) {
	return func(i interface{}) {
		i.(*Adapter).discovery = svc
	}
}

// WithLogger sets a custom logger for this adapter to use.
func WithLogger(logger *log.Logger) func(interface{}) {
	return func(i interface{}) {
		i.(*Adapter).logger = logger
	}
}

// gitHubToken is invoked when specifying a [team.GitHubToken].
func gitHubToken(ctx context.Context, token string) *github.TeamsService {
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return github.NewClient(oauthClient).Teams
}

// discoveryMechanism is invoked when specifying a [team.DiscoveryMechanism].
func discoveryMechanism(ctx context.Context, config map[gosync.ConfigKey]string) (discovery.GitHubDiscovery, error) {
	if config[DiscoveryMechanism] == "saml" {
		oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config[GitHubToken]},
		))

		gitHubV4Client := githubv4.NewClient(oauthClient)
		discoverySvc := saml.New(gitHubV4Client, config[GitHubOrg])

		if strings.ToLower(config[SamlMuteUserNotFoundErr]) == "true" {
			discoverySvc.MuteUserNotFoundErr = true
		}

		return discoverySvc, nil
	}

	return nil, fmt.Errorf("github.team.init -> %w(%s)", gosync.ErrInvalidConfig, config[DiscoveryMechanism])
}

/*
Init a new GitHub Adapter [gosync.Adapter].

Required config:
  - [team.GitHubToken] or [team.WithGitHubV3Client]
  - [team.GitHubOrg]
  - [team.TeamSlug]
  - [team.DiscoveryMechanism] or [team.WithDiscoveryService]

Optional config:
  - [team.SamlMuteUserNotFoundErr]
  - [team.WithLogger]
*/
func Init(
	ctx context.Context,
	config map[gosync.ConfigKey]string,
	optsFn ...func(interface{}),
) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{GitHubOrg, TeamSlug} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("github.team.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	adapter := &Adapter{
		org:    config[GitHubOrg],
		slug:   config[TeamSlug],
		cache:  make(map[string]string),
		logger: log.New(os.Stderr, "[go-sync/github/team] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(adapter)
	}

	if adapter.teams == nil {
		if _, ok := config[GitHubToken]; !ok {
			return nil, fmt.Errorf("github.team.init -> %w(%s)", gosync.ErrMissingConfig, GitHubToken)
		}

		adapter.teams = gitHubToken(ctx, config[GitHubToken])
	}

	if adapter.discovery == nil {
		for _, key := range []gosync.ConfigKey{DiscoveryMechanism, GitHubToken} {
			if _, ok := config[key]; !ok {
				return nil, fmt.Errorf("github.team.init -> %w(%s)", gosync.ErrMissingConfig, key)
			}
		}

		var err error

		adapter.discovery, err = discoveryMechanism(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("github.team.init -> %w", err)
		}
	}

	return adapter, nil
}
