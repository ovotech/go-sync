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

	gosync "github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/types"
)

/*
OpsgenieAPIKey is an API key for authenticating with Opsgenie.
*/
const OpsgenieAPIKey types.ConfigKey = "opsgenie_api_key" //nolint:gosec

// ScheduleID is the name of the Opsgenie Schedule ID.
const ScheduleID types.ConfigKey = "schedule_id"

var (
	_ types.Adapter         = &OnCall{} // Ensure [oncall.OnCall] fully satisfies the [gosync.Adapter] interface.
	_ types.InitFn[*OnCall] = Init      // Ensure [oncall.Init] fully satisfies the [gosync.InitFn] type.
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

// WithClient passes a custom Opsgenie Schedule client to the adapter.
func WithClient(client *schedule.Client) types.ConfigFn[*OnCall] {
	return func(o *OnCall) {
		o.client = client
	}
}

// WithLogger passes a custom logger to the adapter.
func WithLogger(logger *log.Logger) types.ConfigFn[*OnCall] {
	return func(o *OnCall) {
		o.Logger = logger
	}
}

/*
Init a new Opsgenie OnCall [gosync.Adapter].

Required config:
  - [oncall.ScheduleID]
*/
func Init(
	_ context.Context,
	config map[types.ConfigKey]string,
	configFns ...types.ConfigFn[*OnCall],
) (*OnCall, error) {
	for _, key := range []types.ConfigKey{ScheduleID} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("opsgenie.oncall.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	adapter := &OnCall{
		scheduleID: config[ScheduleID],
		getTime:    time.Now,
	}

	if _, ok := config[OpsgenieAPIKey]; ok {
		scheduleClient, err := schedule.NewClient(&client.Config{
			ApiKey: config[OpsgenieAPIKey],
		})
		if err != nil {
			return nil, fmt.Errorf("opsgenie.oncall.init -> %w", err)
		}

		WithClient(scheduleClient)(adapter)
	}

	for _, configFn := range configFns {
		configFn(adapter)
	}

	if adapter.Logger == nil {
		logger := log.New(
			os.Stderr, "[go-sync/opsgenie/oncall] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		)

		WithLogger(logger)(adapter)
	}

	if adapter.client == nil {
		return nil, fmt.Errorf("opsgenie.oncall.init -> %w(%s)", gosync.ErrMissingConfig, OpsgenieAPIKey)
	}

	return adapter, nil
}
