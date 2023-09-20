package schedule_test

import (
	"context"
	"log"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	ogSchedule "github.com/opsgenie/opsgenie-go-sdk-v2/schedule"

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

func ExampleWithClient() {
	ctx := context.Background()

	scheduleClient, err := ogSchedule.NewClient(&client.Config{
		ApiKey: "default",
	})
	if err != nil {
		log.Fatal(err)
	}

	adapter, err := schedule.Init(ctx, map[gosync.ConfigKey]string{
		schedule.ScheduleID: "opsgenie-schedule-id",
	}, schedule.WithClient(scheduleClient))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}

func ExampleWithLogger() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	adapter, err := schedule.Init(ctx, map[gosync.ConfigKey]string{
		schedule.OpsgenieAPIKey: "default",
		schedule.ScheduleID:     "opsgenie-schedule-id",
	}, schedule.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
