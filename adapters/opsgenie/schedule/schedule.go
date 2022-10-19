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
	gosync "github.com/ovotech/go-sync"
	"golang.org/x/exp/slices"
)

/*
OpsgenieAPIKey is an API key for authenticating with Opsgenie.
*/
const OpsgenieAPIKey gosync.ConfigKey = "opsgenie_api_key" //nolint:gosec

// OpsgenieScheduleID is the name of the Opsgenie Schedule ID.
const OpsgenieScheduleID gosync.ConfigKey = "opsgenie_schedule_id"

var (
	_ gosync.Adapter = &Schedule{} // Ensure [schedule.Schedule] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn  = Init        // Ensure the [schedule.Init] function fully satisfies the [gosync.InitFn] type.

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

// New Opsgenie Schedule [gosync.Adapter].
func New(opsgenieConfig *client.Config, scheduleID string, optsFn ...func(schedule *Schedule)) (*Schedule, error) {
	scheduleClient, err := ogSchedule.NewClient(opsgenieConfig)
	if err != nil {
		return nil, fmt.Errorf("opsgenie.schedule.new -> %w", err)
	}

	scheduleAdapter := &Schedule{
		client:     scheduleClient,
		scheduleID: scheduleID,
		Logger:     log.New(os.Stderr, "[gosync/opsgenie/schedule] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(scheduleAdapter)
	}

	return scheduleAdapter, nil
}

/*
Init a new Opsgenie Schedule [gosync.Adapter].

Required config:
  - [schedule.OpsgenieAPIKey]
  - [schedule.OpsgenieScheduleID]
*/
func Init(_ context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{OpsgenieAPIKey, OpsgenieScheduleID} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("opsgenie.oncall.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	opsgenieConfig := client.Config{
		ApiKey: config[OpsgenieAPIKey],
	}

	adapter, err := New(&opsgenieConfig, config[OpsgenieScheduleID])
	if err != nil {
		return nil, fmt.Errorf("opsgenie.oncall.init -> %w", err)
	}

	return adapter, nil
}
