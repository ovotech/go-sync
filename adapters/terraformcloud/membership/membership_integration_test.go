//go:build integration

package membership

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/terraformcloud/membership"
)

var (
	email        = flag.String("email", "test@example.com", "Enter the email of the user for the integration test")
	organisation = flag.String("organisation", "", "Enter the Terraform Cloud organisation name")
	token        = flag.String("token", "", "Enter the Terraform Cloud API token")
)

func TestIntegration(t *testing.T) {
	ctx := context.TODO()

	adapter, err := membership.Init(ctx, map[gosync.ConfigKey]string{
		membership.Token:        *token,
		membership.Organisation: *organisation,
	})
	require.NoError(t, err)

	// Create a membership
	err = adapter.Add(ctx, []string{*email})
	require.NoError(t, err)

	// Assert the membership has been created
	members, err := adapter.Get(ctx)
	require.NoError(t, err)
	assert.Contains(t, members, *email)

	// Delete the membership
	err = adapter.Remove(ctx, []string{*email})
	require.NoError(t, err)

	// Assert the membership has been deleted
	members, err = adapter.Get(ctx)
	require.NoError(t, err)
	assert.NotContains(t, members, *email)
}
