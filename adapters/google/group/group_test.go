package group_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/google/group"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := group.Init(ctx, map[gosync.ConfigKey]string{
		group.GoogleAuthenticationMechanism: "default",
		group.Name:                          "my-group",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
