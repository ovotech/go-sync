package oncall_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/opsgenie/oncall"
)

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
