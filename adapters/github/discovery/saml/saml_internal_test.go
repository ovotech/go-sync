package saml

import (
	"context"
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func emailQueryFunc(names ...string) func(ctx context.Context, q interface{}, variables map[string]interface{}) {
	return func(_ context.Context, query interface{}, variables map[string]interface{}) {
		edges := make([]emailQueryEdge, 0)

		for _, name := range names {
			foo := emailQueryEdge{}
			foo.Node.SamlIdentity.NameID = name
			edges = append(edges, foo)
		}

		arg := query.(*emailQuery)
		arg.Organization.SamlIdentityProvider.ExternalIdentities.Edges = edges
	}
}

func usernameQueryFunc(names ...string) func(ctx context.Context, q interface{}, variables map[string]interface{}) {
	return func(_ context.Context, query interface{}, variables map[string]interface{}) {
		edges := make([]usernameQueryEdge, 0)

		for _, name := range names {
			foo := usernameQueryEdge{}
			foo.Node.User.Login = name
			edges = append(edges, foo)
		}

		arg := query.(*usernameQuery)
		arg.Organization.SamlIdentityProvider.ExternalIdentities.Edges = edges
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	gitHubClient := newMockIGitHubV4Saml(t)
	discovery := New(nil, "test")
	discovery.client = gitHubClient

	assert.Equal(t, "test", discovery.org)
	assert.Zero(t, gitHubClient.Calls)
}

func TestSaml_GetUsernameFromEmail(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		gitHubClient := newMockIGitHubV4Saml(t)
		discovery := New(nil, "test")
		discovery.client = gitHubClient
		discovery.MuteUserNotFoundErr = false

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"email": githubv4.String("foo@email"), "org": githubv4.String("test"),
		},
		).Run(usernameQueryFunc("foo")).Return(nil) //nolint:contextcheck

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"email": githubv4.String("bar@email"), "org": githubv4.String("test"),
		},
		).Run(usernameQueryFunc("bar")).Return(nil) //nolint:contextcheck

		usernames, err := discovery.GetUsernameFromEmail(ctx, []string{"foo@email", "bar@email"})

		require.NoError(t, err)
		assert.ElementsMatch(t, usernames, []string{"foo", "bar"})
	})

	t.Run("MuteUserNotFoundErr false returns error", func(t *testing.T) {
		t.Parallel()

		gitHubClient := newMockIGitHubV4Saml(t)
		discovery := New(nil, "test")
		discovery.client = gitHubClient
		discovery.MuteUserNotFoundErr = false

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"email": githubv4.String("foo@email"), "org": githubv4.String("test"),
		},
		).Run(usernameQueryFunc("foo")).Return(nil) //nolint:contextcheck

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"email": githubv4.String("bar@email"), "org": githubv4.String("test"),
		},
		).Run(usernameQueryFunc()).Return(nil) //nolint:contextcheck

		_, err := discovery.GetUsernameFromEmail(ctx, []string{"foo@email", "bar@email"})

		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("MuteUserNotFoundErr true returns all non-failing results", func(t *testing.T) { //nolint: dupl
		t.Parallel()

		gitHubClient := newMockIGitHubV4Saml(t)
		discovery := New(nil, "test")
		discovery.client = gitHubClient
		discovery.MuteUserNotFoundErr = true

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"email": githubv4.String("foo@email"), "org": githubv4.String("test"),
		},
		).Run(usernameQueryFunc("foo")).Return(nil) //nolint:contextcheck

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"email": githubv4.String("bar@email"), "org": githubv4.String("test"),
		},
		).Run(usernameQueryFunc()).Return(nil) //nolint:contextcheck

		usernames, err := discovery.GetUsernameFromEmail(ctx, []string{"foo@email", "bar@email"})

		require.NoError(t, err)
		assert.ElementsMatch(t, usernames, []string{"foo"})
	})
}

func TestSaml_GetEmailFromUsername(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		gitHubClient := newMockIGitHubV4Saml(t)
		discovery := New(nil, "test")
		discovery.client = gitHubClient

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"login": githubv4.String("foo"), "org": githubv4.String("test"),
		},
		).Run(emailQueryFunc("foo@email")).Return(nil) //nolint:contextcheck

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"login": githubv4.String("bar"), "org": githubv4.String("test"),
		},
		).Run(emailQueryFunc("bar@email")).Return(nil) //nolint:contextcheck

		usernames, err := discovery.GetEmailFromUsername(ctx, []string{"foo", "bar"})

		require.NoError(t, err)
		assert.ElementsMatch(t, usernames, []string{"foo@email", "bar@email"})
	})

	t.Run("MuteUserNotFoundErr false returns error", func(t *testing.T) {
		t.Parallel()

		gitHubClient := newMockIGitHubV4Saml(t)
		discovery := New(nil, "test")
		discovery.client = gitHubClient
		discovery.MuteUserNotFoundErr = false

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"login": githubv4.String("foo"), "org": githubv4.String("test"),
		},
		).Run(emailQueryFunc("foo@email")).Return(nil) //nolint:contextcheck

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"login": githubv4.String("bar"), "org": githubv4.String("test"),
		},
		).Run(emailQueryFunc()).Return(nil) //nolint:contextcheck

		_, err := discovery.GetEmailFromUsername(ctx, []string{"foo", "bar"})

		require.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("MuteUserNotFoundErr true returns all non-failing results", func(t *testing.T) { //nolint: dupl
		t.Parallel()

		gitHubClient := newMockIGitHubV4Saml(t)
		discovery := New(nil, "test")
		discovery.client = gitHubClient
		discovery.MuteUserNotFoundErr = true

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"login": githubv4.String("foo"), "org": githubv4.String("test"),
		},
		).Run(emailQueryFunc("foo@email")).Return(nil) //nolint:contextcheck

		gitHubClient.EXPECT().Query(ctx, mock.Anything, map[string]interface{}{
			"login": githubv4.String("bar"), "org": githubv4.String("test"),
		},
		).Run(emailQueryFunc()).Return(nil) //nolint:contextcheck

		usernames, err := discovery.GetEmailFromUsername(ctx, []string{"foo", "bar"})

		require.NoError(t, err)
		assert.ElementsMatch(t, usernames, []string{"foo@email"})
	})
}
