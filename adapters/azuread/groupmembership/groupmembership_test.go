package groupmembership_test

import (
	"context"
	"log"

	"github.com/ovotech/go-sync/adapters/azuread/groupmembership"
	"github.com/ovotech/go-sync/internal/gosync"
	"github.com/ovotech/go-sync/pkg/types"
)

func ExampleInit() {
	adapter, err := groupmembership.Init(context.TODO(), map[types.ConfigKey]string{
		groupmembership.GroupName: "My Azure AD group",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
