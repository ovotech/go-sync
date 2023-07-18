package schedule_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/opsgenie/schedule"
)

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
