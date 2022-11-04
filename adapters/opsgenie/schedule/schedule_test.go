package schedule_test

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/opsgenie/schedule"
)

func ExampleNew() {
	opsgenieConfig := client.Config{
		ApiKey: "test",
	}

	adapter, err := schedule.New(&opsgenieConfig, "opsgenie-schedule-id")
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleInit() {
	ctx := context.Background()

	adapter, err := schedule.Init(ctx, map[gosync.ConfigKey]string{
		schedule.OpsgenieAPIKey: "default",
		schedule.ScheduleID:     "opsgenie-schedule-id",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
