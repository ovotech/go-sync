package group_test

import (
	"context"
	"log"

	admin "google.golang.org/api/admin/directory/v1"

	"github.com/ovotech/go-sync/packages/google/group"
	"github.com/ovotech/go-sync/packages/gosync"
)

func ExampleNew() {
	ctx := context.Background()

	client, err := admin.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	adapter := group.New(client, "my-group")

	gosync.New(adapter)
}

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
