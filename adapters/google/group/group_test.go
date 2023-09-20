package group_test

import (
	"context"
	"log"
	"os"

	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/google/group"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := group.Init(ctx, map[gosync.ConfigKey]string{
		group.Name: "my-group",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithAdminService() {
	ctx := context.Background()

	client, err := admin.NewService(
		ctx,
		option.WithScopes(admin.AdminDirectoryGroupMemberScope),
		option.WithAPIKey("my-api-key"),
	)
	if err != nil {
		log.Fatal(err)
	}

	adapter, err := group.Init(ctx, map[gosync.ConfigKey]string{
		group.Name: "my-group",
	}, group.WithAdminService(client))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := group.Init(ctx, map[gosync.ConfigKey]string{
		group.Name: "my-group",
	}, group.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
