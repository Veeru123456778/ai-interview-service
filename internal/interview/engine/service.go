package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
)

type Service interface {
	StartInterview(
		ctx context.Context,
		interviewID string,
		userID string,
		resumeID string,
	) (*InterviewState, error)

	GetInterviewState(
		ctx context.Context,
		interviewID string,
	) (*InterviewState, error)

	SaveInterviewState(
		ctx context.Context,
		state *InterviewState,
	) error

	EndInterview(
		ctx context.Context,
		interviewID string,
	) error
}

type service struct {
	repository Repository
	builder    ContextBuilder
}

func NewService(
	repository Repository,
	builder ContextBuilder,
) Service {
	return &service{
		repository: repository,
		builder:    builder,
	}
}

// ----------------------------------------------------------------------
// Start Interview
// ----------------------------------------------------------------------

func (s *service) StartInterview(
	ctx context.Context,
	interviewID string,
	userID string,
	resumeID string,
) (*InterviewState, error) {

	state, err := s.builder.BuildInitialState(
		ctx,
		interviewID,
		userID,
		resumeID,
	)
	if err != nil {
		return nil, fmt.Errorf("build interview state: %w", err)
	}

	if err := s.repository.SaveState(ctx, state); err != nil {
		return nil, fmt.Errorf("save interview state: %w", err)
	}

	return state, nil
}

// ----------------------------------------------------------------------
// Get Interview State
// ----------------------------------------------------------------------

func (s *service) GetInterviewState(
	ctx context.Context,
	interviewID string,
) (*InterviewState, error) {

	state, err := s.repository.GetState(ctx, interviewID)
	if err != nil {
		return nil, fmt.Errorf("get interview state: %w", err)
	}

	return state, nil
}

// ----------------------------------------------------------------------
// Save Interview State
// ----------------------------------------------------------------------

func (s *service) SaveInterviewState(
	ctx context.Context,
	state *InterviewState,
) error {

	state.LastUpdatedAt = time.Now().UTC()

	if err := s.repository.SaveState(ctx, state); err != nil {
		return fmt.Errorf("save interview state: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------
// End Interview
// ----------------------------------------------------------------------

func (s *service) EndInterview(
	ctx context.Context,
	interviewID string,
) error {

	state, err := s.repository.GetState(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("load interview state: %w", err)
	}

	now := time.Now().UTC()

	state.Status = constants.InterviewCompleted
	state.LastUpdatedAt = now
	state.ExpiresAt = now.Add(5 * time.Minute)

	if err := s.repository.SaveState(ctx, state); err != nil {
		return fmt.Errorf("final save interview state: %w", err)
	}

	if err := s.repository.DeleteState(ctx, interviewID); err != nil {
		return fmt.Errorf("delete interview state: %w", err)
	}

	return nil
}