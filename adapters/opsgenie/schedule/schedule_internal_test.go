package schedule

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gosync "github.com/ovotech/go-sync"
)

var errResponse = errors.New("an example error")

func createMockedAdapter(ctx context.Context, t *testing.T) (*Schedule, *mockIOpsgenieSchedule) {
	t.Helper()

	scheduleClient := newMockIOpsgenieSchedule(t)

	adapter, err := Init(ctx, map[gosync.ConfigKey]string{
		OpsgenieAPIKey: "test",
		ScheduleID:     "test",
	})
	require.NoError(t, err)

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
		participants := make([]og.Participant, 0, len(rotationEmails))
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

func TestSchedule_Get(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("error response", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(ctx, t)
		expectedRequest := &schedule.GetRequest{
			IdentifierType:  schedule.Id,
			IdentifierValue: "test",
		}
		scheduleClient.EXPECT().Get(ctx, expectedRequest).Return(nil, errResponse)

		emails, err := adapter.Get(ctx)

		assert.Nil(t, emails)
		require.ErrorContains(t, err, "an example error")
	})

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(ctx, t)
		expectedResponse := testBuildScheduleGetResult(
			1,
			"example1@example.com", "example2@example.com", "example3@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedResponse, nil)

		emails, err := adapter.Get(ctx)

		require.NoError(t, err)
		assert.Equal(t, []string{"example1@example.com", "example2@example.com", "example3@example.com"}, emails)
	})

	t.Run("should handle multiple rotations", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(ctx, t)
		expectedResponse := testBuildScheduleGetResult(
			3,
			"example1@example.com", "example2@example.com", "example3@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedResponse, nil)

		emails, err := adapter.Get(ctx)

		require.NoError(t, err)
		assert.Equal(t, []string{"example1@example.com", "example2@example.com", "example3@example.com"}, emails)
	})

	t.Run("should not duplicate participants across multiple rotations", func(t *testing.T) {
		t.Parallel()

		adapter, scheduleClient := createMockedAdapter(ctx, t)
		expectedResponse := testBuildScheduleGetResult(
			2,
			"example1@example.com", "example2@example.com", "example3@example.com", "example2@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedResponse, nil)

		emails, err := adapter.Get(ctx)

		require.NoError(t, err)
		assert.Equal(t, []string{"example1@example.com", "example2@example.com", "example3@example.com"}, emails)
	})
}

func TestSchedule_Add(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("error response", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(ctx, mock.Anything).Return(nil, errResponse)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		require.ErrorContains(t, err, "an example error")
	})

	t.Run("an error should be returned if no rotations exist", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(0)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		require.ErrorContains(t, err, "gosync cannot create rotations")
	})

	t.Run("an error should be returned if the schedule has more than 1 rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(
			2,
			"example1@example.com", "example2@example.com",
		)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		require.ErrorContains(t, err, "gosync can only manage schedules with a single rotation")
	})

	t.Run("should add new participants to the rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

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

		require.NoError(t, err)
	})

	t.Run("should not add duplicates to the rota", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(
			ctx,
			testBuildExpectedUpdateRotationRequest("example1@example.com", "example2@example.com"),
		).Return(nil, nil)

		err := adapter.Add(ctx, []string{"example2@example.com"})

		require.NoError(t, err)
	})
}

func TestSchedule_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("error response", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(ctx, mock.Anything).Return(nil, errResponse)

		err := adapter.Remove(ctx, []string{"example@example.com"})

		require.ErrorContains(t, err, "an example error")
	})

	t.Run("an error should be returned if no rotations exist", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(0)
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Remove(ctx, []string{})

		require.ErrorContains(t, err, "gosync cannot create rotations")
	})

	t.Run("an error should be returned if the schedule has more than 1 rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(2, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)

		err := adapter.Remove(ctx, []string{"example2@example.com"})

		require.ErrorContains(t, err, "gosync can only manage schedules with a single rotation")
	})

	t.Run("should remove existing participants from the rotation", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(
			ctx,
			testBuildExpectedUpdateRotationRequest("example1@example.com"),
		).Return(nil, nil)

		err := adapter.Remove(ctx, []string{"example2@example.com"})

		require.NoError(t, err)
	})

	t.Run("should ignore nonexistent participants", func(t *testing.T) {
		t.Parallel()
		adapter, scheduleClient := createMockedAdapter(ctx, t)

		expectedScheduleResult := testBuildScheduleGetResult(1, "example1@example.com", "example2@example.com")
		scheduleClient.EXPECT().Get(ctx, mock.Anything).Return(expectedScheduleResult, nil)
		scheduleClient.EXPECT().UpdateRotation(
			ctx,
			testBuildExpectedUpdateRotationRequest("example1@example.com", "example2@example.com"),
		).Return(nil, nil)

		err := adapter.Remove(ctx, []string{"example3@example.com"})

		require.NoError(t, err)
	})
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

		require.NoError(t, err)
		assert.IsType(t, &Schedule{}, adapter)
		assert.Equal(t, "schedule", adapter.scheduleID)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				ScheduleID: "schedule",
			})

			require.ErrorIs(t, err, gosync.ErrMissingConfig)
			require.ErrorContains(t, err, OpsgenieAPIKey)
		})

		t.Run("missing schedule ID", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				OpsgenieAPIKey: "test",
			})

			require.ErrorIs(t, err, gosync.ErrMissingConfig)
			require.ErrorContains(t, err, ScheduleID)
		})
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(os.Stderr, "custom logger", log.LstdFlags)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			OpsgenieAPIKey: "test",
			ScheduleID:     "schedule",
		}, WithLogger(logger))

		require.NoError(t, err)
		assert.Equal(t, logger, adapter.Logger)
	})

	t.Run("with client", func(t *testing.T) {
		t.Parallel()

		scheduleClient, err := schedule.NewClient(&client.Config{
			ApiKey: "test",
		})
		require.NoError(t, err)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			ScheduleID: "schedule",
		}, WithClient(scheduleClient))

		require.NoError(t, err)
		assert.Equal(t, scheduleClient, adapter.client)
	})
}
