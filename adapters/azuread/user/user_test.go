package user_test

import (
	"context"
	"log"

	"github.com/ovotech/go-sync/adapters/azuread/user"
	"github.com/ovotech/go-sync/internal/gosync"
	"github.com/ovotech/go-sync/pkg/types"
)

func ExampleInit() {
	adapter, err := user.Init(context.TODO(), map[types.ConfigKey]string{
		user.Filter: "endsWith(mail, '@example.com')",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
