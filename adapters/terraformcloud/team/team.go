/*
Package team synchronises teams with a Terraform Cloud organisation.

# Requirements

In order to synchronise with Terraform cloud, you will need an Organization API token:
https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/api-tokens#organization-api-tokens

# Examples

See [New] and [Init].
*/
package team

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"

	gosync "github.com/ovotech/go-sync"
)

// Token sets the authentication token for Terraform Cloud.
const Token gosync.ConfigKey = "terraform_cloud_token"

// Organisation sets the Terraform Cloud organisation.
const Organisation gosync.ConfigKey = "terraform_cloud_organisation"

var (
	_ gosync.Adapter = &Team{} // Ensure [team.Team] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn  = Init    // Ensure the [team.Init] function fully satisfies the [gosync.InitFn] type.
)

// iTeams is a subset of Terraform Enterprise Teams, and used to build mocks for easy testing.
type iTeams interface {
	List(ctx context.Context, organization string, options *tfe.TeamListOptions) (*tfe.TeamList, error)
	Create(ctx context.Context, organization string, options tfe.TeamCreateOptions) (*tfe.Team, error)
	Delete(ctx context.Context, teamID string) error
}

type Team struct {
	organisation string
	teams        iTeams
	cache        map[string]string // Cache maps team names to IDs in case they're to be removed.
	Logger       *log.Logger
}

// Get teams in a Terraform Cloud organisation.
func (t *Team) Get(ctx context.Context) ([]string, error) {
	t.Logger.Printf("Fetching teams in Terraform Cloud organisation %s", t.organisation)

	pageNumber := 1
	teams := make([]string, 0)

	t.cache = make(map[string]string)

	t.Logger.Printf("Fetching first page")

	for {
		tfeTeams, err := t.teams.List(ctx, t.organisation, &tfe.TeamListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber},
		})
		if err != nil {
			return nil, fmt.Errorf("teams.list(%s) -> %w", t.organisation, err)
		}

		t.Logger.Printf("Fetched page %v in %v", tfeTeams.CurrentPage, tfeTeams.TotalPages)

		for _, team := range tfeTeams.Items {
			teams = append(teams, team.Name)
			t.cache[team.Name] = team.ID
		}

		pageNumber = tfeTeams.NextPage

		if tfeTeams.CurrentPage >= tfeTeams.TotalPages {
			break
		}
	}

	t.Logger.Println("Fetched teams successfully")

	return teams, nil
}

// Add teams to a Terraform Cloud organisation.
func (t *Team) Add(ctx context.Context, teams []string) error {
	t.Logger.Printf("Adding %s to Terraform Cloud organisation %s", teams, t.organisation)

	for _, team := range teams {
		team := team

		_, err := t.teams.Create(ctx, t.organisation, tfe.TeamCreateOptions{Name: &team})
		if err != nil {
			return fmt.Errorf("terraformcloud.team.add -> %w", err)
		}
	}

	t.Logger.Println("Finished adding teams successfully")

	return nil
}

// Remove teams from a Terraform Cloud organisation.
func (t *Team) Remove(ctx context.Context, teams []string) error {
	t.Logger.Printf("Removing %s from Terraform Cloud organisation %s", teams, t.organisation)

	for _, team := range teams {
		err := t.teams.Delete(ctx, t.cache[team])
		if err != nil {
			return fmt.Errorf("terraformcloud.team.remove -> %w", err)
		}
	}

	t.Logger.Println("Finished removing teams successfully")

	return nil
}

// WithClient passes a custom Terraform Cloud client to the adapter.
func WithClient(client *tfe.Client) gosync.ConfigFn {
	return func(i interface{}) {
		if adapter, ok := i.(*Team); ok {
			adapter.teams = client.Teams
		}
	}
}

// WithLogger passes a custom logger to the adapter.
func WithLogger(logger *log.Logger) gosync.ConfigFn {
	return func(i interface{}) {
		if adapter, ok := i.(*Team); ok {
			adapter.Logger = logger
		}
	}
}

/*
Init a new Terraform Cloud Team [gosync.Adapter].

Required config:
  - [team.Organisation]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string, configFns ...gosync.ConfigFn) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{Organisation} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("team.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	adapter := &Team{
		organisation: config[Organisation],
		cache:        make(map[string]string),
	}

	if _, ok := config[Token]; ok {
		client, err := tfe.NewClient(&tfe.Config{Token: config[Token]})
		if err != nil {
			return nil, fmt.Errorf("team.init.newclient -> %w", err)
		}

		WithClient(client)(adapter)
	}

	for _, configFn := range configFns {
		configFn(adapter)
	}

	if adapter.Logger == nil {
		logger := log.New(
			os.Stderr, "[go-sync/terraformcloud/team] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		)

		WithLogger(logger)(adapter)
	}

	if adapter.teams == nil {
		return nil, fmt.Errorf("team.init -> %w(%s)", gosync.ErrMissingConfig, Token)
	}

	return adapter, nil
}
