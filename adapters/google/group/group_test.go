package group_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/google/group"
	admin "google.golang.org/api/admin/directory/v1"
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
		group.GoogleGroupName:               "my-group",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
