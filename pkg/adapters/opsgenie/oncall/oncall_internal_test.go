package oncall

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/ovotech/go-sync/internal/mocks"
	"github.com/stretchr/testify/assert"
)

var errGetOnCall = errors.New("an example error")

func createMockedAdapter(t *testing.T, clock clock) (*OnCall, *mocks.IOpsgenieSchedule) {
	t.Helper()

	scheduleClient := mocks.NewIOpsgenieSchedule(t)
	adapter, _ := New(&client.Config{
		ApiKey: "test",
	}, "test")
	adapter.client = scheduleClient
	adapter.clock = clock

	return adapter, scheduleClient
}

func TestNew(t *testing.T) {
	t.Parallel()

	scheduleClient := mocks.NewIOpsgenieSchedule(t)
	adapter, err := New(&client.Config{
		ApiKey: "test",
	}, "test")
	adapter.client = scheduleClient

	assert.NoError(t, err)
	assert.Equal(t, "test", adapter.scheduleID)
	assert.Zero(t, scheduleClient.Calls)
}

func TestOnCall_Get(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expectedTime := time.Date(2022, 10, 6, 12, 0, 0, 0, time.UTC)
	flat := true

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		clock := mocks.NewClock(t)
		clock.On("Now").Return(expectedTime)
		adapter, scheduleClient := createMockedAdapter(t, clock)
		expectedRequest := &schedule.GetOnCallsRequest{
			Flat:                   &flat,
			Date:                   &expectedTime,
			ScheduleIdentifierType: schedule.Id,
			ScheduleIdentifier:     "test",
		}
		scheduleClient.EXPECT().GetOnCalls(ctx, expectedRequest).Return(&schedule.GetOnCallsResult{
			OnCallRecipients: []string{"foo@email.com", "bar@email.com"},
		}, nil)

		emails, err := adapter.Get(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"foo@email.com", "bar@email.com"}, emails)
	})

	t.Run("error response", func(t *testing.T) {
		t.Parallel()

		clock := mocks.NewClock(t)
		clock.On("Now").Return(expectedTime)
		adapter, scheduleClient := createMockedAdapter(t, clock)
		expectedRequest := &schedule.GetOnCallsRequest{
			Flat:                   &flat,
			Date:                   &expectedTime,
			ScheduleIdentifierType: schedule.Id,
			ScheduleIdentifier:     "test",
		}
		scheduleClient.EXPECT().GetOnCalls(ctx, expectedRequest).Return(nil, errGetOnCall)

		emails, err := adapter.Get(ctx)

		assert.Nil(t, emails)
		assert.Error(t, err, "an example error")
	})
}

func TestOnCall_Add(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clock := mocks.NewClock(t)
	adapter, scheduleClient := createMockedAdapter(t, clock)

	err := adapter.Add(ctx, []string{"example@bar.com"})

	assert.Error(t, err, "not implemented")
	assert.Zero(t, scheduleClient.Calls)
}

func TestOnCall_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clock := mocks.NewClock(t)
	adapter, scheduleClient := createMockedAdapter(t, clock)

	err := adapter.Remove(ctx, []string{"example@bar.com"})

	assert.Error(t, err, "not implemented")
	assert.Zero(t, scheduleClient.Calls)
}
