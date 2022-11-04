/*
Package oncall allows you to synchronise other services with the emails of users who are currently on-call for a
schedule.

Note: On-call is readonly, and so you can only use this as a source.

# Requirements

You will need to create an [API Key] with the following access rights:
  - Read

# Examples

See [New] and [Init].

[API Key]: https://support.atlassian.com/opsgenie/docs/api-key-management/
*/
package oncall

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	gosync "github.com/ovotech/go-sync"
)

/*
OpsgenieAPIKey is an API key for authenticating with Opsgenie.
*/
const OpsgenieAPIKey gosync.ConfigKey = "opsgenie_api_key" //nolint:gosec

// ScheduleID is the name of the Opsgenie Schedule ID.
const ScheduleID gosync.ConfigKey = "schedule_id"

var (
	_ gosync.Adapter = &OnCall{} // Ensure [oncall.OnCall] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn  = Init      // Ensure the [oncall.Init] function fully satisfies the [gosync.InitFn] type.
)

type iOpsgenieSchedule interface {
	GetOnCalls(context context.Context, request *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error)
}

type OnCall struct {
	client     iOpsgenieSchedule
	scheduleID string
	getTime    func() time.Time
	Logger     *log.Logger
}

// Get email addresses of users currently on-call.
func (o *OnCall) Get(ctx context.Context) ([]string, error) {
	o.Logger.Printf("Fetching users currently on-call in Opsgenie schedule %s", o.scheduleID)

	date := o.getTime()
	flat := true
	onCallRequest := &schedule.GetOnCallsRequest{
		Flat:                   &flat,
		Date:                   &date,
		ScheduleIdentifierType: schedule.Id,
		ScheduleIdentifier:     o.scheduleID,
	}

	result, err := o.client.GetOnCalls(ctx, onCallRequest)
	if err != nil {
		return nil, fmt.Errorf("opsgenie.oncall.get.getoncalls -> %w", err)
	}

	o.Logger.Println("Fetched on-call users successfully")

	return result.OnCallRecipients, nil
}

// Add is not supported, as the on-call is readonly.
func (o *OnCall) Add(_ context.Context, _ []string) error {
	return gosync.ErrReadOnly
}

// Remove is not supported, as the on-call is readonly.
func (o *OnCall) Remove(_ context.Context, _ []string) error {
	return gosync.ErrReadOnly
}

// New Opsgenie OnCall [gosync.Adapter].
func New(opsgenieConfig *client.Config, scheduleID string, optsFn ...func(schedule *OnCall)) (*OnCall, error) {
	scheduleClient, err := schedule.NewClient(opsgenieConfig)
	if err != nil {
		return nil, fmt.Errorf("opsgenie.oncall.new -> %w", err)
	}

	onCallAdapter := &OnCall{
		client:     scheduleClient,
		scheduleID: scheduleID,
		getTime:    time.Now,
		Logger:     log.New(os.Stderr, "[go-sync/opsgenie/oncall]", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(onCallAdapter)
	}

	return onCallAdapter, nil
}

/*
Init a new Opsgenie OnCall [gosync.Adapter].

Required config:
  - [oncall.OpsgenieAPIKey]
  - [oncall.ScheduleID]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{OpsgenieAPIKey, ScheduleID} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("opsgenie.oncall.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	opsgenieConfig := client.Config{
		ApiKey: config[OpsgenieAPIKey],
	}

	adapter, err := New(&opsgenieConfig, config[ScheduleID])
	if err != nil {
		return nil, fmt.Errorf("opsgenie.oncall.init -> %w", err)
	}

	return adapter, nil
}
