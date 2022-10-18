# Opsgenie Schedule adapter for Go Sync

This adapter allows you to synchronise the participants of a schedule. Using this as a source supports schedule with
multiple rotations, however if you wish to use this as a destination adapter the schedule must only have 1 rotation
configured, and all members of the source adapter must already have an Opsgenie license allocated.

## Requirements

You will need to create an [API Key](https://support.atlassian.com/opsgenie/docs/api-key-management/) with the following
permissions:

| Access rights |
|:--------------|
| Read          |
| Update        |

## Example

```go
package main

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/opsgenie/schedule"
)

func main() {
	opsgenieConfig := client.Config{
		ApiKey: "test",
	}
	scheduleAdapter, err := schedule.New(&opsgenieConfig, "opsgenie-schedule-id")

	svc := gosync.New(scheduleAdapter)

	// Synchronise the participants of a schedule with something else.
	anotherServiceAdapter := someAdapter.New()

	err = svc.SyncWith(context.Background(), anotherServiceAdapter)
	if err != nil {
		log.Fatal(err)
	}
}
```
