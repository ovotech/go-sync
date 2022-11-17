/*
Package team synchronises teams with Terraform Cloud.
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

// iTeams is a subset of the Terraform Enterprise Teams Client, and used to build mocks for easy testing.
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

	tfeTeams, err := t.teams.List(ctx, t.organisation, &tfe.TeamListOptions{})
	if err != nil {
		return nil, fmt.Errorf("terraformcloud.team.get.list -> %w", err)
	}

	teams := make([]string, 0, len(tfeTeams.Items))

	for _, team := range tfeTeams.Items {
		teams = append(teams, team.Name)
		t.cache[team.Name] = team.ID
	}

	t.Logger.Println("Fetched teams successfully")

	return teams, nil
}

// Add things to Team service.
func (t *Team) Add(ctx context.Context, teams []string) error {
	t.Logger.Printf("Adding %s to Terraform Cloud organisation %s", teams, t.organisation)

	for _, team := range teams {
		_, err := t.teams.Create(ctx, t.organisation, tfe.TeamCreateOptions{Name: &team})
		if err != nil {
			return fmt.Errorf("terraformcloud.team.add -> %w", err)
		}
	}

	t.Logger.Println("Finished adding teams successfully")

	return nil
}

// Remove things from Team service.
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

// New Team [gosync.adapter].
func New(client tfe.Client, organisation string) *Team {
	return &Team{
		teams:        client.Teams,
		organisation: organisation,
		cache:        make(map[string]string),
		Logger:       log.New(os.Stderr, "[go-sync/terraformcloud/team] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}
}

/*
Init a new Terraform Cloud Team [gosync.Adapter].

Required config:
  - [team.Token]
  - [team.Organisation]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{Token, Organisation} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("team.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	client, err := tfe.NewClient(&tfe.Config{Token: config[Token]})
	if err != nil {
		return nil, fmt.Errorf("team.init.newclient -> %w", err)
	}

	return New(*client, config[Organisation]), nil
}
