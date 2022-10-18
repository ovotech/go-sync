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

// Ensure the adapter type fully satisfies the gosync.Adapter interface.
var _ gosync.Adapter = &Schedule{}

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
	logger     *log.Logger
}

var ErrMultipleRotations = errors.New("gosync can only manage schedules with a single rotation")
var ErrNoRotations = errors.New("gosync cannot create rotations - you must have 1 already defined for schedule")

// New instantiates a new Opsgenie Schedule adapter.
func New(opsgenieConfig *client.Config, scheduleID string, optsFn ...func(schedule *Schedule)) (*Schedule, error) {
	scheduleClient, err := ogSchedule.NewClient(opsgenieConfig)
	if err != nil {
		return nil, fmt.Errorf("opsgenie.schedule.new -> %w", err)
	}

	scheduleAdapter := &Schedule{
		client:     scheduleClient,
		scheduleID: scheduleID,
		logger:     log.New(os.Stderr, "[gosync/opsgenie/schedule] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, fn := range optsFn {
		fn(scheduleAdapter)
	}

	return scheduleAdapter, nil
}

// Get fetches a flattened list of all participants, even across multiple rotations.
func (s *Schedule) Get(ctx context.Context) ([]string, error) {
	s.logger.Printf("Getting all participants in the Opsgenie schedule %s", s.scheduleID)

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

	s.logger.Printf("Found %d participants of schedule %s: %s", len(emails), s.scheduleID, emails)

	return emails, nil
}

// Add lets you add new participants to a rotation, but the schedule must only have 1 rotation defined.
func (s *Schedule) Add(ctx context.Context, emails []string) error {
	s.logger.Printf("Adding %d users to schedule %s: %s", len(emails), s.scheduleID, emails)

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

// Remove lets you remove participants from a rotation, but the schedule must only have 1 rotation defined.
func (s *Schedule) Remove(ctx context.Context, emails []string) error {
	s.logger.Printf("Removing %d users from schedule %s: %s", len(emails), s.scheduleID, emails)

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
		s.logger.Printf("Fetching schedule %s from Opsgenie", s.scheduleID)

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
		s.logger.Printf("Already have schedule %s cached", s.scheduleID)
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

	s.logger.Printf("Updating rotation %s of schedule %s with participants: %s", rotation.Id, s.scheduleID, emails)

	_, err := s.client.UpdateRotation(ctx, request)
	if err != nil {
		return fmt.Errorf("error when updating participants: %w", err)
	}

	return nil
}
