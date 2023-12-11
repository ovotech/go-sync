//go:build integration

package usergroup

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	gosync "github.com/ovotech/go-sync"
)

var group = flag.String("group", "", "Enter the user group ID to adjust group membership of in the Integration test")
var email = flag.String("email", "", "Enter the email of a user to add to the user group membership of in the Integration test")
var key = os.Getenv("SLACK_API_KEY")

func TestIntegration(t *testing.T) {
	if *group == "" {
		t.Fatalf("Required parameter 'group' is missing or empty.")
	}

	if *email == "" {
		t.Fatalf("Required parameter 'email' is missing or empty.")
	}

	if key == "" {
		t.Fatalf("Required environment variable 'SLACK_API_KEY' is missing or empty.")
	}

	ctx := context.TODO()
	adapter, err := Init(ctx, map[gosync.ConfigKey]string{
		UserGroupID: *group,
		SlackAPIKey: key,
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	emails, err := adapter.Get(ctx)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.NotContainsf(t, emails, *email, "Email %s already exists in the user group %s", *email, *group)

	err = adapter.Add(ctx, []string{*email})
	assert.NoError(t, err)

	err = adapter.Remove(ctx, emails)
	assert.NoError(t, err)

	newemails, err := adapter.Get(ctx)
	assert.NoError(t, err)
	assert.Containsf(t, newemails, *email, "Email %s not found after adding it to the group", *email)

	err = adapter.Add(ctx, emails)
	assert.NoError(t, err)

	err = adapter.Remove(ctx, []string{*email})
	assert.NoError(t, err)

	finalemails, err := adapter.Get(ctx)
	assert.NoError(t, err)
	assert.NotContainsf(t, finalemails, *email, "Email %s still exists after removing it from the group %s", *email, *group)
}
