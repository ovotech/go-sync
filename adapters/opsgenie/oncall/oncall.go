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

// Ensure the Opsgenie OnCall adapter type fully satisfies the gosync.Adapter interface.
var _ gosync.Adapter = &OnCall{}

const (
	/*
		API key for authenticating with Opsgenie.
	*/
	OpsgenieAPIKey           gosync.ConfigKey = "opsgenie_api_key"            //nolint:gosec
	OpsgenieOnCallScheduleID gosync.ConfigKey = "opsgenie_oncall_schedule_id" // Schedule ID.
)

type iOpsgenieSchedule interface {
	GetOnCalls(context context.Context, request *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error)
}

type OnCall struct {
	client     iOpsgenieSchedule
	scheduleID string
	getTime    func() time.Time
	logger     *log.Logger
}

// New instantiates a new Opsgenie OnCall adapter.
func New(opsgenieConfig *client.Config, scheduleID string, optsFn ...func(schedule *OnCall)) (*OnCall, error) {
	scheduleClient, err := schedule.NewClient(opsgenieConfig)

	if err != nil {
		return nil, fmt.Errorf("opsgenie.oncall.new -> %w", err)
	}

	onCallAdapter := &OnCall{
		client:     scheduleClient,
		scheduleID: scheduleID,
		getTime:    time.Now,
		logger:     log.New(os.Stderr, "[go-sync/opsgenie/oncall]", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(onCallAdapter)
	}

	return onCallAdapter, nil
}

// Ensure the Init function fully satisfies the gosync.InitFn type.
var _ gosync.InitFn = Init

// Init a new Opsgenie OnCall gosync.Adapter. All gosync.ConfigKey keys are required in config.
func Init(config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{OpsgenieAPIKey, OpsgenieOnCallScheduleID} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("opsgenie.oncall.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	opsgenieConfig := client.Config{
		ApiKey: config[OpsgenieAPIKey],
	}

	adapter, err := New(&opsgenieConfig, config[OpsgenieOnCallScheduleID])
	if err != nil {
		return nil, fmt.Errorf("opsgenie.oncall.init -> %w", err)
	}

	return adapter, nil
}

// Get emails of users currently on-call in on-call.
func (o *OnCall) Get(ctx context.Context) ([]string, error) {
	o.logger.Printf("Fetching users currently on-call in Opsgenie schedule %s", o.scheduleID)

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

	o.logger.Println("Fetched on-call users successfully")

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
