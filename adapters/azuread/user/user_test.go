package user_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/azuread/user"
)

func ExampleInit() {
	adapter, err := user.Init(context.TODO(), map[gosync.ConfigKey]string{
		user.Filter: "endsWith(mail, '@example.com')",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
