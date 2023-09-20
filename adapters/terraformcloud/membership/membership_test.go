package membership_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/terraformcloud/membership"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := membership.Init(ctx, map[gosync.ConfigKey]string{
		membership.Token:        "my-org-token",
		membership.Organisation: "ovotech",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
