package membership_test

import (
	"context"
	"log"

	"github.com/ovotech/go-sync/adapters/terraformcloud/membership"
	"github.com/ovotech/go-sync/internal/gosync"
	"github.com/ovotech/go-sync/pkg/types"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := membership.Init(ctx, map[types.ConfigKey]string{
		membership.Token:        "my-org-token",
		membership.Organisation: "ovotech",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
