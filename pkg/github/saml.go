package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/shurcooL/githubv4"
)

// ErrUserNotFound is returned when the SAML -> GitHub username mapping cannot be found.
var ErrUserNotFound = errors.New("user_not_found")

type Saml struct {
	client *githubv4.Client
	org    string
}

// NewSamlDiscoveryService instantiates a new GitHub SAML discovery service.
func NewSamlDiscoveryService(client *githubv4.Client, org string) *Saml {
	return &Saml{
		org:    org,    // GitHub organisation.
		client: client, // GitHub V4 GraphQL client.
	}
}

// GetEmailFromUsername takes a list of GitHub usernames, and returns a list of emails.
func (s *Saml) GetEmailFromUsername(logins ...string) ([]string, error) { //nolint:dupl
	emails := make([]string, 0, len(logins))

	for _, login := range logins {
		var query struct {
			Organization struct {
				SamlIdentityProvider struct {
					ExternalIdentities struct {
						Edges []struct {
							Node struct {
								SamlIdentity struct {
									NameId string //nolint:revive,stylecheck
								}
							}
						}
					} `graphql:"externalIdentities(login: $login, first: 1)"`
				}
			} `graphql:"organization(login: $org)"`
		}

		variables := map[string]interface{}{
			"org":   githubv4.String(s.org),
			"login": githubv4.String(login),
		}

		if err := s.client.Query(context.Background(), &query, variables); err != nil {
			return nil, fmt.Errorf("github.saml.GetEmailFromUsername(%s) -> graphql query error -> %w", login, err)
		}

		if len(query.Organization.SamlIdentityProvider.ExternalIdentities.Edges) != 1 {
			return nil, fmt.Errorf("github.saml.GetEmailFromUsername(%s) -> unknown identity -> %w", login, ErrUserNotFound)
		}

		emails = append(
			emails,
			query.Organization.SamlIdentityProvider.ExternalIdentities.Edges[0].Node.SamlIdentity.NameId,
		)
	}

	return emails, nil
}

// GetUsernameFromEmail takes a list of email addresses, and returns a list of equivalent GitHub usernames.
func (s *Saml) GetUsernameFromEmail(emails ...string) ([]string, error) { //nolint:dupl
	ids := make([]string, 0, len(emails))

	for _, email := range emails {
		var query struct {
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

		variables := map[string]interface{}{
			"org":   githubv4.String(s.org),
			"email": githubv4.String(email),
		}

		if err := s.client.Query(context.Background(), &query, variables); err != nil {
			return nil, fmt.Errorf("github.saml.GetUsernameFromEmail(%s) -> graphql query error -> %w", email, err)
		}

		if len(query.Organization.SamlIdentityProvider.ExternalIdentities.Edges) != 1 {
			return nil, fmt.Errorf("github.saml.GetUsernameFromEmail(%s) -> %w", email, ErrUserNotFound)
		}

		ids = append(ids, query.Organization.SamlIdentityProvider.ExternalIdentities.Edges[0].Node.User.Login)
	}

	return ids, nil
}
