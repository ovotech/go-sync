package user_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/terraformcloud/user"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := user.Init(ctx, map[gosync.ConfigKey]string{
		user.Token:        "my-org-token",
		user.Organisation: "ovotech",
		user.Team:         "my-team",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
