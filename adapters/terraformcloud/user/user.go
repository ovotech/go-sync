/*
Package user synchronises users in a Terraform Cloud team.

# Requirements

In order to synchronise with Terraform cloud, you will need an Organization API token:
https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/api-tokens#organization-api-tokens

# Examples

See [New] and [Init].
*/
package user

import (
	"context"
	"errors"
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

// Team is the name of the team to sync with.
const Team gosync.ConfigKey = "team"

var (
	_ gosync.Adapter = &User{} // Ensure [user.User] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn  = Init    // Ensure the [user.Init] function fully satisfies the [gosync.InitFn] type.
)

// ErrTeamNotFound is returned if the team cannot be found in the Terraform Cloud organisation.
var ErrTeamNotFound = errors.New("team_not_found")

// iTeams is a subset of Terraform Enterprise TeamMembers, and used to build mocks for easy testing.
type iTeamMembers interface {
	List(ctx context.Context, teamID string) ([]*tfe.User, error)
	Add(ctx context.Context, teamID string, options tfe.TeamMemberAddOptions) error
	Remove(ctx context.Context, teamID string, options tfe.TeamMemberRemoveOptions) error
}

// iTeams is a subset of Terraform Enterprise Teams, and used to build mocks for easy testing.
type iTeams interface {
	List(ctx context.Context, organization string, options *tfe.TeamListOptions) (*tfe.TeamList, error)
}

type iOrganizationMemberships interface {
	List(
		ctx context.Context,
		organization string,
		options *tfe.OrganizationMembershipListOptions,
	) (*tfe.OrganizationMembershipList, error)
}

type User struct {
	organisation            string
	team                    string
	teams                   iTeams
	teamMembers             iTeamMembers
	organizationMemberships iOrganizationMemberships
	Logger                  *log.Logger
}

// getTeamID queries the Terraform Cloud API to convert a friendly team name into a team ID.
func (u *User) getTeamID(ctx context.Context) (string, error) {
	u.Logger.Printf("Querying Terraform Cloud organisation %s for team ID of %s", u.organisation, u.team)

	teams, err := u.teams.List(ctx, u.organisation, &tfe.TeamListOptions{Names: []string{u.team}})
	if err != nil {
		return "", fmt.Errorf("terraformcloud.user.get(%s, %s) -> %w", u.organisation, u.team, err)
	}

	if len(teams.Items) != 1 {
		return "", fmt.Errorf("terraformcloud.user.get(%s, %s) -> %w", u.organisation, u.team, ErrTeamNotFound)
	}

	u.Logger.Println("Successfully queried team ID")

	return teams.Items[0].ID, nil
}

// getOrgIDsFromEmails takes a slice of emails, and returns a slice of Organisational Membership IDs.
func (u *User) getOrgIDsFromEmails(ctx context.Context, emails []string) ([]string, error) {
	pageNumber := 1
	ids := make([]string, 0, len(emails))

	u.Logger.Printf("Fetching %s from Terraform Cloud organisation %s", emails, u.organisation)

	for {
		users, err := u.organizationMemberships.List(ctx, u.organisation, &tfe.OrganizationMembershipListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber},
			Emails:      emails,
		})
		if err != nil {
			return nil, fmt.Errorf("organizationmembership.list(%s, %s) -> %w", u.organisation, u.team, err)
		}

		u.Logger.Printf("Fetching page %v in %v", users.CurrentPage, users.TotalPages)

		for _, user := range users.Items {
			ids = append(ids, user.ID)
		}

		pageNumber = users.NextPage

		if users.CurrentPage >= users.TotalPages {
			break
		}
	}

	u.Logger.Println("Finished fetching users")

	return ids, nil
}

// Get users in a Terraform Cloud team.
func (u *User) Get(ctx context.Context) ([]string, error) {
	u.Logger.Printf("Fetching users in Terraform Cloud team %s", u.team)

	team, err := u.teams.List(ctx, u.organisation, &tfe.TeamListOptions{
		Include: []tfe.TeamIncludeOpt{tfe.TeamOrganizationMemberships},
		Names:   []string{u.team},
	})
	if err != nil {
		return nil, fmt.Errorf("terraformcloud.user.get(%s).list(%s) -> %w", u.organisation, u.team, err)
	}

	if len(team.Items) != 1 {
		return nil, fmt.Errorf("terraformcloud.user.get(%s).list(%s) -> %w", u.organisation, u.team, ErrTeamNotFound)
	}

	emails := make([]string, 0, len(team.Items[0].OrganizationMemberships))

	for _, organisationMembership := range team.Items[0].OrganizationMemberships {
		emails = append(emails, organisationMembership.Email)
	}

	u.Logger.Println("Fetched teams successfully")

	return emails, nil
}

// Add users to a Terraform Cloud team.
func (u *User) Add(ctx context.Context, emails []string) error {
	u.Logger.Printf("Adding %s to Terraform Cloud team %s", emails, u.team)

	ids, err := u.getOrgIDsFromEmails(ctx, emails)
	if err != nil {
		return fmt.Errorf(
			"terraformcloud.user.add(%s, %s).getorgidsfromemails(%s) -> %w",
			u.organisation, u.team, emails, err,
		)
	}

	teamID, err := u.getTeamID(ctx)
	if err != nil {
		return fmt.Errorf("terraformcloud.user.add(%s, %s).getteamid -> %w", u.organisation, u.team, err)
	}

	err = u.teamMembers.Add(ctx, teamID, tfe.TeamMemberAddOptions{OrganizationMembershipIDs: ids})
	if err != nil {
		return fmt.Errorf("terraformcloud.user.add(%s, %s).add(%s) -> %w", u.organisation, u.team, emails, err)
	}

	u.Logger.Println("Finished adding users successfully")

	return nil
}

// Remove users from a Terraform Cloud team.
func (u *User) Remove(ctx context.Context, emails []string) error {
	u.Logger.Printf("Removing %s from Terraform Cloud team %s", emails, u.team)

	ids, err := u.getOrgIDsFromEmails(ctx, emails)
	if err != nil {
		return fmt.Errorf(
			"terraformcloud.user.remove(%s, %s).getorgidsfromemails(%s) -> %w",
			u.organisation, u.team, emails, err,
		)
	}

	teamID, err := u.getTeamID(ctx)
	if err != nil {
		return fmt.Errorf("terraformcloud.user.remove(%s, %s).getteamid -> %w", u.organisation, u.team, err)
	}

	err = u.teamMembers.Remove(ctx, teamID, tfe.TeamMemberRemoveOptions{OrganizationMembershipIDs: ids})
	if err != nil {
		return fmt.Errorf("terraformcloud.user.remove(%s, %s).add(%s) -> %w", u.organisation, u.team, emails, err)
	}

	u.Logger.Println("Finished removing users successfully")

	return nil
}

// New Terraform Cloud User [gosync.adapter].
func New(client *tfe.Client, organisation string, team string) *User {
	return &User{
		organisation:            organisation,
		team:                    team,
		teams:                   client.Teams,
		teamMembers:             client.TeamMembers,
		organizationMemberships: client.OrganizationMemberships,
		Logger: log.New(
			os.Stderr, "[go-sync/terraformcloud/user] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		),
	}
}

/*
Init a new Terraform Cloud User [gosync.Adapter].

Required config:
  - [user.Token]
  - [user.Team]
  - [user.Organisation]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{Token, Organisation, Team} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("team.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	client, err := tfe.NewClient(&tfe.Config{Token: config[Token]})
	if err != nil {
		return nil, fmt.Errorf("team.init.newclient -> %w", err)
	}

	return New(client, config[Organisation], config[Team]), nil
}
