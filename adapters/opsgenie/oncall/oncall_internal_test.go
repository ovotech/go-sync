package oncall

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/stretchr/testify/assert"

	gosync "github.com/ovotech/go-sync"
)

var errGetOnCall = errors.New("an example error")

func createMockedAdapter(t *testing.T, mockedTime time.Time) (*OnCall, *mockIOpsgenieSchedule) {
	t.Helper()

	scheduleClient := newMockIOpsgenieSchedule(t)
	adapter, _ := New(&client.Config{
		ApiKey: "test",
	}, "test")
	adapter.client = scheduleClient
	adapter.getTime = func() time.Time {
		return mockedTime
	}

	return adapter, scheduleClient
}

func TestNew(t *testing.T) {
	t.Parallel()

	scheduleClient := newMockIOpsgenieSchedule(t)
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

		adapter, scheduleClient := createMockedAdapter(t, expectedTime)
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

		adapter, scheduleClient := createMockedAdapter(t, expectedTime)
		expectedRequest := &schedule.GetOnCallsRequest{
			Flat:                   &flat,
			Date:                   &expectedTime,
			ScheduleIdentifierType: schedule.Id,
			ScheduleIdentifier:     "test",
		}
		scheduleClient.EXPECT().GetOnCalls(ctx, expectedRequest).Return(nil, errGetOnCall)

		emails, err := adapter.Get(ctx)

		assert.Nil(t, emails)
		assert.ErrorContains(t, err, "an example error")
	})
}

func TestOnCall_Add(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	adapter, scheduleClient := createMockedAdapter(t, time.Now())

	err := adapter.Add(ctx, []string{"example@bar.com"})

	assert.ErrorIs(t, err, gosync.ErrReadOnly)
	assert.Zero(t, scheduleClient.Calls)
}

func TestOnCall_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	adapter, scheduleClient := createMockedAdapter(t, time.Now())

	err := adapter.Remove(ctx, []string{"example@bar.com"})

	assert.ErrorIs(t, err, gosync.ErrReadOnly)
	assert.Zero(t, scheduleClient.Calls)
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			OpsgenieAPIKey: "test",
			ScheduleID:     "schedule",
		})

		assert.NoError(t, err)
		assert.IsType(t, &OnCall{}, adapter)
		assert.Equal(t, "schedule", adapter.(*OnCall).scheduleID)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				ScheduleID: "schedule",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, OpsgenieAPIKey)
		})

		t.Run("missing schedule ID", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				OpsgenieAPIKey: "test",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, ScheduleID)
		})
	})
}
