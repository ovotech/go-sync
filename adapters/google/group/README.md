# Google Groups adapter for Go Sync
This adapter synchronises email addresses with a Google Group.

## Requirements
In order to synchronise with Google, you'll need to credentials with the Admin SDK enabled on your account, and 
credentials with the following scopes:

| Go Scopes                              | Google Scopes                                                |
|----------------------------------------|--------------------------------------------------------------|
| `admin.AdminDirectoryGroupMemberScope` | https://www.googleapis.com/auth/admin.directory.group.member |


## Example
```go
package main

import (
	"context"
	"log"

	"github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/google/group"
	admin "google.golang.org/api/admin/directory/v1"
)

func main() {
	ctx := context.Background()

	client, err := admin.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	googleGroup := group.New(client, "my-group")

	svc := gosync.New(googleGroup)

	// Synchronise a Google Group with something else.
	anotherServiceAdapter := someAdapter.New()

	err = svc.SyncWith(context.Background(), anotherServiceAdapter)
	if err != nil {
		log.Fatal(err)
	}
}
```
