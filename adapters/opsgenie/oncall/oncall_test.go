package oncall_test

import (
	"context"
	"log"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"

	"github.com/ovotech/go-sync/adapters/opsgenie/oncall"
	"github.com/ovotech/go-sync/internal/gosync"
	"github.com/ovotech/go-sync/pkg/types"
)

func ExampleInit() {
	ctx := context.Background()

	adapter, err := oncall.Init(ctx, map[types.ConfigKey]string{
		oncall.OpsgenieAPIKey: "default",
		oncall.ScheduleID:     "opsgenie-schedule-id",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithClient() {
	ctx := context.Background()

	scheduleClient, err := schedule.NewClient(&client.Config{
		ApiKey: "default",
	})
	if err != nil {
		log.Fatal(err)
	}

	adapter, err := oncall.Init(ctx, map[types.ConfigKey]string{
		oncall.ScheduleID: "opsgenie-schedule-id",
	}, oncall.WithClient(scheduleClient))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := oncall.Init(ctx, map[types.ConfigKey]string{
		oncall.OpsgenieAPIKey: "default",
		oncall.ScheduleID:     "opsgenie-schedule-id",
	}, oncall.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
