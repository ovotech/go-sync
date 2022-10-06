# Opsgenie On-Call adapter for Go Sync

This adapter allows you to synchronise other services with the emails of users who are currently on-call for a schedule.

**Note:** On-call is readonly, and so you can only use this as a source.

## Requirements

You will need to create an [API Key](https://support.atlassian.com/opsgenie/docs/api-key-management/) with the following
permissions:

| Access rights |
|:--------------|
| Read          |

## Example

```go
package main

import (
	"context"
	"log"
	
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/ovotech/go-sync/pkg/adapters/opsgenie/oncall"
	"github.com/ovotech/go-sync/pkg/sync"
)

func main() {
	opsgenieConfig := client.Config{
		ApiKey: "test",
	}
	onCallAdapter := oncall.New(&opsgenieConfig, "opsgenie-schedule-id")

	svc := sync.New(onCallAdapter)

	// Synchronise an on-call list with something else.
	anotherServiceAdapter := someAdapter.New()
	
	err := svc.SyncWith(context.Background(), anotherServiceAdapter)
	if err != nil {
		log.Fatal(err)
	}
}
```
