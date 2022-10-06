package oncall

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/ovotech/go-sync/internal/types"
	"github.com/ovotech/go-sync/pkg/ports"
)

// Ensure the adapter type fully satisfies the ports.Adapter interface.
var _ ports.Adapter = &OnCall{}

// ErrNotImplemented is returned if this adapter is used as a destination.
var ErrNotImplemented = errors.New("not implemented - on-call is readonly")

type iOpsgenieSchedule interface {
	GetOnCalls(context context.Context, request *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error)
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

type OnCall struct {
	client     iOpsgenieSchedule
	scheduleID string
	clock      types.Clock
	logger     types.Logger
}

// New instantiates a new Opsgenie OnCall adapter.
func New(opsgenieConfig *client.Config, scheduleID string, optsFn ...func(schedule *OnCall)) *OnCall {
	logger := log.New(os.Stderr, "[go-sync/opsgenie/oncall]", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	scheduleClient, err := schedule.NewClient(opsgenieConfig)

	if err != nil {
		logger.Fatalf("Error occurred when creating on-call client: %s", err)

		return nil
	}

	onCallAdapter := &OnCall{
		client:     scheduleClient,
		scheduleID: scheduleID,
		clock:      realClock{},
		logger:     logger,
	}

	for _, fn := range optsFn {
		fn(onCallAdapter)
	}

	return onCallAdapter
}

// Get emails of users currently on-call in on-call.
func (o *OnCall) Get(ctx context.Context) ([]string, error) {
	o.logger.Printf("Fetching users currently on-call in Opsgenie schedule %o", o.scheduleID)

	date := o.clock.Now()
	flat := true
	onCallRequest := &schedule.GetOnCallsRequest{
		Flat:                   &flat,
		Date:                   &date,
		ScheduleIdentifierType: schedule.Id,
		ScheduleIdentifier:     o.scheduleID,
	}

	result, err := o.client.GetOnCalls(ctx, onCallRequest)
	if err != nil {
		return nil, fmt.Errorf("opsgenie.oncall.get.GetOnCalls -> %w", err)
	}

	o.logger.Println("Fetched on-call users successfully")

	return result.OnCallRecipients, nil
}

// Add is not supported, as the on-call is readonly.
func (o *OnCall) Add(_ context.Context, _ []string) error {
	return ErrNotImplemented
}

// Remove is not supported, as the on-call is readonly.
func (o *OnCall) Remove(_ context.Context, _ []string) error {
	return ErrNotImplemented
}
