package groupmembership_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/azuread/groupmembership"
)

func ExampleInit() {
	adapter, err := groupmembership.Init(context.TODO(), map[gosync.ConfigKey]string{
		groupmembership.GroupName: "My Azure AD group",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
