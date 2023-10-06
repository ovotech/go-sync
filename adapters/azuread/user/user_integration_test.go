//go:build integration

package user

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	gosync "github.com/ovotech/go-sync/pkg/types"
)

var email = flag.String("email", "test@example.com", "Enter the email of a user to search for in the Integration test")

func TestIntegration(t *testing.T) {
	ctx := context.TODO()
	adapter, err := Init(ctx, map[gosync.ConfigKey]string{
		Filter: fmt.Sprintf("mail eq '%s'", *email),
	})
	assert.NoError(t, err)

	emails, err := adapter.Get(ctx)
	assert.NoError(t, err)
	assert.ElementsMatch(t, emails, []string{*email})
}
