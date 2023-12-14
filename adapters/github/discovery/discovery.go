package discovery

import "context"

// GitHubDiscovery is required because there are multiple ways to convert a GitHub email into a username.
// At OVO we use SAML, but other organisations may use public emails or another mechanism.
type GitHubDiscovery interface {
	GetUsernameFromEmail(ctx context.Context, emails []string) ([]string, error)
	GetEmailFromUsername(ctx context.Context, usernames []string) ([]string, error)
}
