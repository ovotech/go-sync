package saml

import (
	"context"
	"errors"
	"fmt"

	"github.com/ovotech/go-sync/adapters/github/discovery"
	"github.com/shurcooL/githubv4"
)

// Ensure the Saml adapter type fully satisfies the discovery.GitHubDiscovery interface.
var _ discovery.GitHubDiscovery = &Saml{}

type iGitHubV4Saml interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

type emailQuery struct {
	Organization struct {
		SamlIdentityProvider struct {
			ExternalIdentities struct {
				Edges []struct {
					Node struct {
						SamlIdentity struct {
							NameID string
						}
					}
				}
			} `graphql:"externalIdentities(login: $login, first: 1)"`
		}
	} `graphql:"organization(login: $org)"`
}

type usernameQuery struct {
	Organization struct {
		SamlIdentityProvider struct {
			ExternalIdentities struct {
				Edges []struct {
					Node struct {
						User struct {
							Login string
						}
					}
				}
			} `graphql:"externalIdentities(userName: $email, first: 1)"`
		}
	} `graphql:"organization(login: $org)"`
}

// ErrUserNotFound is returned when the SAML -> GitHub username mapping cannot be found.
var ErrUserNotFound = errors.New("user_not_found")

type Saml struct {
	client iGitHubV4Saml
	org    string // GitHub organisation.
}

// New instantiates a new GitHub SAML discovery adapter for use with GitHub adapters.
func New(client *githubv4.Client, org string, optsFn ...func(*Saml)) *Saml {
	saml := &Saml{
		client: client,
		org:    org,
	}

	for _, fn := range optsFn {
		fn(saml)
	}

	return saml
}

// GetEmailFromUsername takes a list of GitHub usernames, and returns a list of emails.
func (s *Saml) GetEmailFromUsername(ctx context.Context, logins []string) ([]string, error) {
	emails := make([]string, 0, len(logins))

	for _, login := range logins {
		query := &emailQuery{}

		variables := map[string]interface{}{
			"org":   githubv4.String(s.org),
			"login": githubv4.String(login),
		}

		err := s.client.Query(ctx, query, variables)
		if err != nil {
			return nil, fmt.Errorf("github.saml.getemailfromusername(%s) -> graphql query error -> %w", login, err)
		}

		if len(query.Organization.SamlIdentityProvider.ExternalIdentities.Edges) != 1 {
			return nil, fmt.Errorf("github.saml.getemailfromusername(%s) -> unknown identity -> %w", login, ErrUserNotFound)
		}

		emails = append(
			emails,
			query.Organization.SamlIdentityProvider.ExternalIdentities.Edges[0].Node.SamlIdentity.NameID,
		)
	}

	return emails, nil
}

// GetUsernameFromEmail takes a list of email addresses, and returns a list of equivalent GitHub usernames.
func (s *Saml) GetUsernameFromEmail(ctx context.Context, emails []string) ([]string, error) {
	ids := make([]string, 0, len(emails))

	for _, email := range emails {
		query := &usernameQuery{}

		variables := map[string]interface{}{
			"org":   githubv4.String(s.org),
			"email": githubv4.String(email),
		}

		if err := s.client.Query(ctx, query, variables); err != nil {
			return nil, fmt.Errorf("github.saml.getusernamefromemail(%s) -> graphql query error -> %w", email, err)
		}

		if len(query.Organization.SamlIdentityProvider.ExternalIdentities.Edges) != 1 {
			return nil, fmt.Errorf("github.saml.getusernamefromemail(%s) -> %w", email, ErrUserNotFound)
		}

		ids = append(ids, query.Organization.SamlIdentityProvider.ExternalIdentities.Edges[0].Node.User.Login)
	}

	return ids, nil
}
