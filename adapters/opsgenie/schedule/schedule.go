/*
Package schedule allows you to synchronise the participants of a schedule.

Using this as a source supports schedule with multiple rotations, however if you wish to use this as a destination
adapter the schedule must only have 1 rotation configured, and all members of the source adapter must already have an
Opsgenie license allocated.

# Requirements

You will need to create an [API Key] with the following access rights:
  - Read
  - Update

# Examples

See [New] and [Init].

[API Key]: https://support.atlassian.com/opsgenie/docs/api-key-management/
*/
package schedule

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	ogSchedule "github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"golang.org/x/exp/slices"

	gosync "github.com/ovotech/go-sync"
)

/*
OpsgenieAPIKey is an API key for authenticating with Opsgenie.
*/
const OpsgenieAPIKey gosync.ConfigKey = "opsgenie_api_key" //nolint:gosec

// ScheduleID is the name of the Opsgenie Schedule ID.
const ScheduleID gosync.ConfigKey = "schedule_id"

var (
	_ gosync.Adapter           = &Schedule{} // Ensure [schedule.Schedule] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn[*Schedule] = Init        // Ensure [schedule.Init] fully satisfies the [gosync.InitFn] type.

	ErrMultipleRotations = errors.New("gosync can only manage schedules with a single rotation")
	ErrNoRotations       = errors.New("gosync cannot create rotations - you must have 1 already defined for schedule")
)

type iOpsgenieSchedule interface {
	Get(ctx context.Context, request *ogSchedule.GetRequest) (*ogSchedule.GetResult, error)
	UpdateRotation(
		ctx context.Context,
		request *ogSchedule.UpdateRotationRequest,
	) (*ogSchedule.UpdateRotationResult, error)
}

type Schedule struct {
	client     iOpsgenieSchedule
	scheduleID string
	schedule   *ogSchedule.GetResult
	Logger     *log.Logger
}

// Get a flattened list of all participants, even across multiple rotations.
func (s *Schedule) Get(ctx context.Context) ([]string, error) {
	s.Logger.Printf("Getting all participants in the Opsgenie schedule %s", s.scheduleID)

	result, err := s.fetchSchedule(ctx)
	if err != nil {
		return nil, fmt.Errorf("opsgenie.schedule.get.fetchschedule -> %w", err)
	}

	emails := make([]string, 0)

	for _, rotation := range result.Schedule.Rotations {
		for _, participant := range rotation.Participants {
			if participant.Type == og.User && !slices.Contains(emails, participant.Username) {
				emails = append(emails, participant.Username)
			}
		}
	}

	s.Logger.Printf("Found %d participants of schedule %s: %s", len(emails), s.scheduleID, emails)

	return emails, nil
}

// Add new participants to a rotation, but the schedule must only have 1 rotation defined.
func (s *Schedule) Add(ctx context.Context, emails []string) error {
	s.Logger.Printf("Adding %d users to schedule %s: %s", len(emails), s.scheduleID, emails)

	result, err := s.fetchSchedule(ctx)
	if err != nil {
		return fmt.Errorf("opsgenie.schedule.add.fetchschedule -> %w", err)
	}

	rotation, err := s.getRotation(result)
	if err != nil {
		return fmt.Errorf("opsgenie.schedule.add.getrotation -> %w", err)
	}

	// Get current participants
	updatedParticipants := make([]string, 0)

	for _, participant := range rotation.Participants {
		if participant.Type == og.User {
			updatedParticipants = append(updatedParticipants, participant.Username)
		}
	}

	// Add any new participants
	for _, email := range emails {
		if !slices.Contains(updatedParticipants, email) {
			updatedParticipants = append(updatedParticipants, email)
		}
	}

	err = s.updateParticipants(ctx, rotation, updatedParticipants)
	if err != nil {
		return fmt.Errorf("opsgenie.schedule.add.updateparticipants -> %w", err)
	}

	return nil
}

// Remove participants from a rotation, but the schedule must only have 1 rotation defined.
func (s *Schedule) Remove(ctx context.Context, emails []string) error {
	s.Logger.Printf("Removing %d users from schedule %s: %s", len(emails), s.scheduleID, emails)

	result, err := s.fetchSchedule(ctx)
	if err != nil {
		return fmt.Errorf("opsgenie.schedule.remove.fetchschedule -> %w", err)
	}

	rotation, err := s.getRotation(result)
	if err != nil {
		return fmt.Errorf("opsgenie.schedule.remove.getrotation -> %w", err)
	}

	// Build up the participants list - removing any that are to be removed
	updatedParticipants := make([]string, 0)

	for _, participant := range rotation.Participants {
		if participant.Type == og.User && !slices.Contains(emails, participant.Username) {
			updatedParticipants = append(updatedParticipants, participant.Username)
		}
	}

	err = s.updateParticipants(ctx, rotation, updatedParticipants)
	if err != nil {
		return fmt.Errorf("opsgenie.schedule.remove.updateparticipants -> %w", err)
	}

	return nil
}

func (s *Schedule) getRotation(result *ogSchedule.GetResult) (*og.Rotation, error) {
	if len(result.Schedule.Rotations) > 1 {
		return nil, fmt.Errorf("getrotation(%s) -> %w", s.scheduleID, ErrMultipleRotations)
	}

	if len(result.Schedule.Rotations) == 0 {
		return nil, fmt.Errorf("getrotations(%s) -> %w", s.scheduleID, ErrNoRotations)
	}

	return &result.Schedule.Rotations[0], nil
}

func (s *Schedule) fetchSchedule(ctx context.Context) (*ogSchedule.GetResult, error) {
	if s.schedule == nil {
		s.Logger.Printf("Fetching schedule %s from Opsgenie", s.scheduleID)

		scheduleRequest := &ogSchedule.GetRequest{
			IdentifierType:  ogSchedule.Id,
			IdentifierValue: s.scheduleID,
		}

		result, err := s.client.Get(ctx, scheduleRequest)
		if err != nil {
			return nil, fmt.Errorf("error when fetching schedules: %w", err)
		}

		s.schedule = result
	} else {
		s.Logger.Printf("Already have schedule %s cached", s.scheduleID)
	}

	return s.schedule, nil
}

func (s *Schedule) updateParticipants(ctx context.Context, rotation *og.Rotation, emails []string) error {
	participants := make([]og.Participant, 0)
	for _, email := range emails {
		participants = append(participants, og.Participant{
			Type:     og.User,
			Username: email,
		})
	}

	request := &ogSchedule.UpdateRotationRequest{
		ScheduleIdentifierType:  ogSchedule.Id,
		ScheduleIdentifierValue: s.scheduleID,
		RotationId:              rotation.Id,
		Rotation: &og.Rotation{
			Participants: participants,
		},
	}

	s.Logger.Printf("Updating rotation %s of schedule %s with participants: %s", rotation.Id, s.scheduleID, emails)

	_, err := s.client.UpdateRotation(ctx, request)
	if err != nil {
		return fmt.Errorf("error when updating participants: %w", err)
	}

	return nil
}

// WithClient passes a custom Opsgenie Schedule client to the adapter.
func WithClient(client *ogSchedule.Client) gosync.ConfigFn[*Schedule] {
	return func(s *Schedule) {
		s.client = client
	}
}

// WithLogger passes a custom logger to the adapter.
func WithLogger(logger *log.Logger) gosync.ConfigFn[*Schedule] {
	return func(s *Schedule) {
		s.Logger = logger
	}
}

/*
Init a new Opsgenie Schedule [gosync.Adapter].

Required config:
  - [schedule.ScheduleID]
*/
func Init(
	_ context.Context,
	config map[gosync.ConfigKey]string,
	configFns ...gosync.ConfigFn[*Schedule],
) (*Schedule, error) {
	for _, key := range []gosync.ConfigKey{ScheduleID} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("opsgenie.schedule.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	adapter := &Schedule{
		scheduleID: config[ScheduleID],
	}

	if _, ok := config[OpsgenieAPIKey]; ok {
		scheduleClient, err := ogSchedule.NewClient(&client.Config{
			ApiKey: config[OpsgenieAPIKey],
		})
		if err != nil {
			return nil, fmt.Errorf("opsgenie.schedule.init -> %w", err)
		}

		WithClient(scheduleClient)(adapter)
	}

	for _, configFn := range configFns {
		configFn(adapter)
	}

	if adapter.Logger == nil {
		logger := log.New(
			os.Stderr, "[go-sync/opsgenie/schedule] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		)

		WithLogger(logger)(adapter)
	}

	if adapter.client == nil {
		return nil, fmt.Errorf("opsgenie.schedule.init -> %w(%s)", gosync.ErrMissingConfig, OpsgenieAPIKey)
	}

	return adapter, nil
}
