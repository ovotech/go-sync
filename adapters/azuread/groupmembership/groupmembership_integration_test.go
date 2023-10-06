//go:build integration

package groupmembership

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	gosync "github.com/ovotech/go-sync/pkg/types"
)

var group = flag.String("group", "", "Enter the group to adjust group membership of in the Integration test")

var email = flag.String("email", "", "Enter the email of a user to change membership of in the Integration test")

func TestIntegration(t *testing.T) {
	if *group == "" {
		t.Fatalf("Required parameter 'group' is missing or empty.")
	}

	if *email == "" {
		t.Fatalf("Required parameter 'email' is missing or empty.")
	}

	ctx := context.TODO()
	adapter, err := Init(ctx, map[gosync.ConfigKey]string{
		GroupName: *group,
	})
	assert.NoError(t, err)

	emails, err := adapter.Get(ctx)
	assert.NoError(t, err)
	assert.NotContainsf(t, emails, *email, "Email %s already exists in the group %s", *email, *group)

	err = adapter.Add(ctx, []string{*email})
	assert.NoError(t, err)

	emails, err = adapter.Get(ctx)
	assert.NoError(t, err)
	assert.Containsf(t, emails, *email, "Email %s not found after adding it to the group", *email)

	err = adapter.Remove(ctx, []string{*email})
	assert.NoError(t, err)

	emails, err = adapter.Get(ctx)
	assert.NoError(t, err)
	assert.NotContainsf(t, emails, *email, "Email %s still exists after removing it from the group %s", *email, *group)
}
