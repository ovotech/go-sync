package schedule

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	gosync "github.com/ovotech/go-sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errResponse = errors.New("an example error")

func createMockedAdapter(t *testing.T) (*Schedule, *mockIOpsgenieSchedule) {
	t.Helper()

	scheduleClient := newMockIOpsgenieSchedule(t)
	adapter, _ := New(&client.Config{
		ApiKey: "test",
	}, "test")
	adapter.client = scheduleClient

	return adapter, scheduleClient
}

func testBuildExpectedUpdateRotationRequest(emails ...string) *schedule.UpdateRotationRequest {
	participants := make([]og.Participant, len(emails))
	for i, email := range emails {
		participants[i] = og.Participant{
			Type:     og.User,
			Username: email,
		}
	}

	return &schedule.UpdateRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: "test",
		RotationId:              "rotation-0",
		Rotation: &og.Rotation{
			Participants: participants,
		},
	}
}

func testBuildScheduleGetResult(numRotations int, emails ...string) *schedule.GetResult {
	var chunkedEmails [][]string

	if numRotations > 0 {
		chunkSize := len(emails) / numRotations
		for chunkSize < len(emails) {
			emails, chunkedEmails = emails[chunkSize:], append(chunkedEmails, emails[0:chunkSize:chunkSize])
		}

		chunkedEmails = append(chunkedEmails, emails)
	}

	rotations := make([]og.Rotation, numRotations)

	for index, rotationEmails := range chunkedEmails {
		var participants = make([]og.Participant, len(rotationEmails))
		for _, email := range rotationEmails {
			participants = append(participants, og.Participant{
				Type:     og.User,
				Username: email,
			})
		}

		rotations[index] = og.Rotation{
			Id:              fmt.Sprintf("rotation-%d", index),
			Name:            fmt.Sprintf("Example Rotation %d", index),
			StartDate:       nil,
			EndDate:         nil,
			Type:            og.Weekly,
			Length:          uint32(len(participants)),
			Participants:    participants,
			TimeRestriction: nil,
		}
	}

	return &schedule.GetResult{
		Schedule: schedule.Schedule{
			Id:          "test",
			Name:        "Example",
			Description: "",
			Timezone:    "UTC",
			Enabled:     true,
			OwnerTeam: &og.OwnerTeam{
				Id:   "team",
				Name: "Team Name",
			},
			Rotations: rotations,
		},
	}
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

func TestSchedule_Get(t *testing.T) { //nolint:funlen
	t.Parallel()

	ctx := context.Background()

	t.Run("error response", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(t)
		expectedRequest := &schedule.GetRequest{
			IdentifierType:  schedule.Id,
			IdentifierValue: "test",
		}
		scheduleClient.EXPECT().Get(ctx, expectedRequest).Return(nil, errResponse)

		emails, err := adapter.Get(ctx)

		assert.Nil(t, emails)
		assert.ErrorContains(t, err, "an example error")
	})

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(t)
		expectedResponse := testBuildScheduleGetResult(
			1,
			"example1@example.com", "example2@example.com", "example3@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedResponse, nil)

		emails, err := adapter.Get(ctx)

		assert.Nil(t, err)
		assert.Equal(t, emails, []string{"example1@example.com", "example2@example.com", "example3@example.com"})
	})

	t.Run("should handle multiple rotations", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(t)
		expectedResponse := testBuildScheduleGetResult(
			3,
			"example1@example.com", "example2@example.com", "example3@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedResponse, nil)

		emails, err := adapter.Get(ctx)

		assert.Nil(t, err)
		assert.Equal(t, emails, []string{"example1@example.com", "example2@example.com", "example3@example.com"})
	})

	t.Run("should not duplicate participants across multiple rotations", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(t)
		expectedResponse := testBuildScheduleGetResult(
			2,
			"example1@example.com", "example2@example.com", "example3@example.com", "example2@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedResponse, nil)

		emails, err := adapter.Get(ctx)

		assert.Nil(t, err)
		assert.Equal(t, emails, []string{"example1@example.com", "example2@example.com", "example3@example.com"})
	})
}

func TestSchedule_Add(t *testing.T) { //nolint:funlen
	t.Parallel()

	ctx := context.Background()

	t.Run("error response", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(ctx, mock.Anything).Return(nil, errResponse)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		assert.ErrorContains(t, err, "an example error")
	})

	t.Run("an error should be returned if no rotations exist", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(0)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		assert.ErrorContains(t, err, "gosync cannot create rotations")
	})

	t.Run("an error should be returned if the schedule has more than 1 rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(
			2,
			"example1@example.com", "example2@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		assert.ErrorContains(t, err, "gosync can only manage schedules with a single rotation")
	})

	t.Run("should add new participants to the rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(
			1,
			"example1@example.com", "example2@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(
			ctx,
			testBuildExpectedUpdateRotationRequest("example1@example.com", "example2@example.com", "example3@example.com"),
		).Return(nil, nil)

		err := adapter.Add(ctx, []string{"example3@example.com"})

		assert.Nil(t, err)
	})

	t.Run("should not add duplicates to the rota", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(
			ctx,
			testBuildExpectedUpdateRotationRequest("example1@example.com", "example2@example.com"),
		).Return(nil, nil)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		assert.Nil(t, err)
	})
}

func TestSchedule_Remove(t *testing.T) { //nolint:funlen
	t.Parallel()

	ctx := context.Background()

	t.Run("error response", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(ctx, mock.Anything).Return(nil, errResponse)

		err := adapter.Remove(ctx, []string{"example@example.com"})

		assert.ErrorContains(t, err, "an example error")
	})

	t.Run("an error should be returned if no rotations exist", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(0)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Remove(ctx, []string{})

		assert.ErrorContains(t, err, "gosync cannot create rotations")
	})

	t.Run("an error should be returned if the schedule has more than 1 rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(2, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Remove(ctx, []string{"example2@example.com"})

		assert.ErrorContains(t, err, "gosync can only manage schedules with a single rotation")
	})

	t.Run("should remove existing participants from the rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(
			ctx,
			testBuildExpectedUpdateRotationRequest("example1@example.com"),
		).Return(nil, nil)

		err := adapter.Remove(ctx, []string{"example2@example.com"})

		assert.Nil(t, err)
	})

	t.Run("should ignore nonexistent participants", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(
			ctx,
			testBuildExpectedUpdateRotationRequest("example1@example.com", "example2@example.com"),
		).Return(nil, nil)

		err := adapter.Remove(ctx, []string{"example3@example.com"})

		assert.Nil(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(map[gosync.ConfigKey]string{
			OpsgenieAPIKey:     "test",
			OpsgenieScheduleID: "schedule",
		})

		assert.NoError(t, err)
		assert.IsType(t, &Schedule{}, adapter)
		assert.Equal(t, "schedule", adapter.(*Schedule).scheduleID)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[gosync.ConfigKey]string{
				OpsgenieScheduleID: "schedule",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, OpsgenieAPIKey)
		})

		t.Run("missing schedule ID", func(t *testing.T) {
			t.Parallel()

			_, err := Init(map[gosync.ConfigKey]string{
				OpsgenieAPIKey: "test",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, OpsgenieScheduleID)
		})
	})
}
