/*
Package membership synchronises members in a Terraform Cloud organisation.

# Requirements

In order to synchronise with Terraform cloud, you will need an Organization API token:
https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/api-tokens#organization-api-tokens

# Examples

See [New] and [Init].
*/
package membership

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
	_ gosync.Adapter             = &Membership{} // Ensure [team.Team] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn[*Membership] = Init          // Ensure [team.Init] fully satisfies the [gosync.InitFn] type.
)

// iOrganizationMemberships is a subset of Terraform Enterprise
// OrganizationMemberships, and used to build mocks for easy testing.
type iOrganizationMemberships interface {
	List(ctx context.Context, organization string, options *tfe.OrganizationMembershipListOptions) (
		*tfe.OrganizationMembershipList, error)
	Create(ctx context.Context, organization string, options tfe.OrganizationMembershipCreateOptions) (
		*tfe.OrganizationMembership, error)
	Read(ctx context.Context, organizationMembershipID string) (*tfe.OrganizationMembership, error)
	Delete(ctx context.Context, organizationMembershipID string) error
}

type Membership struct {
	organisation            string
	organizationMemberships iOrganizationMemberships
	Logger                  *log.Logger
}

// getOrgIDsFromEmails takes a slice of emails, and returns a slice of Organisational Membership IDs.
func (m *Membership) getOrgIDsFromEmails(ctx context.Context, emails []string) ([]string, error) {
	pageNumber := 1
	ids := make([]string, 0, len(emails))

	m.Logger.Printf("Fetching IDs from Terraform Cloud organisation %s", m.organisation)

	for {
		users, err := m.organizationMemberships.List(ctx, m.organisation, &tfe.OrganizationMembershipListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber},
			Emails:      emails,
		})
		if err != nil {
			return nil, fmt.Errorf("terraformcloud.membership.getOrgIDsFromEmails(%s).list -> %w", emails, err)
		}

		m.Logger.Printf("Fetching page %v in %v", users.CurrentPage, users.TotalPages)

		for _, user := range users.Items {
			ids = append(ids, user.ID)
		}

		pageNumber = users.NextPage

		if users.CurrentPage >= users.TotalPages {
			break
		}
	}

	m.Logger.Println("Finished fetching users")

	return ids, nil
}

// Get memberships in a Terraform Cloud organisation.
func (m *Membership) Get(ctx context.Context) ([]string, error) {
	pageNumber := 1
	memberships := make([]string, 0)

	m.Logger.Printf("Fetching members in Terraform Cloud organisation %s", m.organisation)

	for {
		listOptions := &tfe.OrganizationMembershipListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: pageNumber,
			},
		}

		tfeMemberships, err := m.organizationMemberships.List(ctx, m.organisation, listOptions)
		if err != nil {
			return nil, fmt.Errorf("terraformcloud.membership.get.list(%s) -> %w", m.organisation, err)
		}

		for _, membership := range tfeMemberships.Items {
			memberships = append(memberships, membership.Email)
		}

		pageNumber = tfeMemberships.NextPage

		if tfeMemberships.CurrentPage >= tfeMemberships.TotalPages {
			break
		}

		m.Logger.Printf("Fetching page %v in %v", tfeMemberships.CurrentPage, tfeMemberships.TotalPages)
	}

	m.Logger.Println("Fetched memberships successfully")

	return memberships, nil
}

// Add members to a Terraform Cloud organisation.
func (m *Membership) Add(ctx context.Context, emails []string) error {
	m.Logger.Printf("Adding %s to Terraform Cloud organisation %s", emails, m.organisation)

	for _, email := range emails {
		options := tfe.OrganizationMembershipCreateOptions{
			Email: &email,
			Type:  "organization-memberships",
		}

		_, err := m.organizationMemberships.Create(ctx, m.organisation, options)
		if err != nil {
			return fmt.Errorf("terraformcloud.membership.add(%s).create(%s) -> %w", emails, email, err)
		}
	}

	m.Logger.Println("Finished adding members successfully")

	return nil
}

// Remove members from the Terraform Cloud organisation.
func (m *Membership) Remove(ctx context.Context, emails []string) error {
	m.Logger.Printf("Removing %s from Terraform Cloud organisation %s", emails, m.organisation)

	ids, err := m.getOrgIDsFromEmails(ctx, emails)
	if err != nil {
		return fmt.Errorf(
			"terraformcloud.membership.remove(%s).getorgidsfromemails -> %w",
			emails, err,
		)
	}

	for _, id := range ids {
		err = m.organizationMemberships.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("terraformcloud.membership.remove(%s).delete(%s) -> %w", emails, id, err)
		}
	}

	m.Logger.Println("Finished removing members successfully")

	return nil
}

// WithClient passes a custom Terraform Cloud client to the adapter.
func WithClient(client *tfe.Client) gosync.ConfigFn[*Membership] {
	return func(u *Membership) {
		u.organizationMemberships = client.OrganizationMemberships
	}
}

// WithLogger passes a custom logger to the adapter.
func WithLogger(logger *log.Logger) gosync.ConfigFn[*Membership] {
	return func(u *Membership) {
		u.Logger = logger
	}
}

/*
Init a new Terraform Cloud Membership [gosync.Adapter].

Required config:
  - [membership.Token]
  - [membership.Organisation]
*/
func Init(
	_ context.Context,
	config map[gosync.ConfigKey]string,
	configFns ...gosync.ConfigFn[*Membership],
) (*Membership, error) {
	for _, key := range []gosync.ConfigKey{Organisation} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("terraformcloud.membership.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	adapter := &Membership{
		organisation: config[Organisation],
	}

	if _, ok := config[Token]; ok {
		client, err := tfe.NewClient(&tfe.Config{Token: config[Token]})
		if err != nil {
			return nil, fmt.Errorf("user.init.newclient -> %w", err)
		}

		WithClient(client)(adapter)
	}

	for _, configFn := range configFns {
		configFn(adapter)
	}

	if adapter.Logger == nil {
		logger := log.New(
			os.Stderr, "[go-sync/terraformcloud/membership] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		)

		WithLogger(logger)(adapter)
	}

	return adapter, nil
}
