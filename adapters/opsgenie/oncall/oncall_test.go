package oncall_test

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/opsgenie/oncall"
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
		oncall.OpsgenieAPIKey:     "default",
		oncall.OpsgenieScheduleID: "opsgenie-schedule-id",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
