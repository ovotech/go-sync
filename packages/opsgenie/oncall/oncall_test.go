package oncall_test

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"

	"github.com/ovotech/go-sync/packages/gosync"
	"github.com/ovotech/go-sync/packages/opsgenie/oncall"
)

func ExampleNew() {
	opsgenieConfig := client.Config{
		ApiKey: "test",
	}

	adapter, err := oncall.New(&opsgenieConfig, "opsgenie-schedule-id")
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleInit() {
	ctx := context.Background()

	adapter, err := oncall.Init(ctx, map[gosync.ConfigKey]string{
		oncall.OpsgenieAPIKey: "default",
		oncall.ScheduleID:     "opsgenie-schedule-id",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
