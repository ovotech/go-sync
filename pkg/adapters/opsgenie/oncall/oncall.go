package oncall

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/ovotech/go-sync/internal/types"
	gosyncerrors "github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/ports"
)

// Ensure the adapter type fully satisfies the ports.Adapter interface.
var _ ports.Adapter = &OnCall{}

type iOpsgenieSchedule interface {
	GetOnCalls(context context.Context, request *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error)
}

type OnCall struct {
	client     iOpsgenieSchedule
	scheduleID string
	getTime    func() time.Time
	logger     types.Logger
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

// Get emails of users currently on-call in on-call.
func (o *OnCall) Get(ctx context.Context) ([]string, error) {
	o.logger.Printf("Fetching users currently on-call in Opsgenie schedule %o", o.scheduleID)

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
	return gosyncerrors.ErrReadOnly
}

// Remove is not supported, as the on-call is readonly.
func (o *OnCall) Remove(_ context.Context, _ []string) error {
	return gosyncerrors.ErrReadOnly
}
