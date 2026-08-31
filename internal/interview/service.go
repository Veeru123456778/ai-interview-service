package interview

import (
	"context"
	"fmt"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
	"github.com/Veeru123456778/ai-interview-service/internal/resume"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
	"github.com/google/uuid"
)

// ----------------------------------------------------------------------
// Service Interface
// ----------------------------------------------------------------------

type Service interface {
	CreateInterview(
		ctx context.Context,
		userID string,
		request *CreateInterviewRequest,
	) (*CreateInterviewResponse, error)

	GetInterview(
		ctx context.Context,
		interviewID string,
		userID string,
	) (*InterviewResponse, error)

	ListInterviews(
		ctx context.Context,
		userID string,
	) (*ListInterviewsResponse, error)

	StartInterview(
		ctx context.Context,
		interviewID string,
		userID string,
	) (*engine.InterviewState, error)

	GetInterviewState(
		ctx context.Context,
		interviewID string,
	) (*engine.InterviewState, error)

	UpdateInterviewState(
		ctx context.Context,
		state *engine.InterviewState,
	) error

	CompleteInterview(
		ctx context.Context,
		interviewID string,
		userID string,
	) error
}

// ----------------------------------------------------------------------
// Service Implementation
// ----------------------------------------------------------------------

type service struct {
	repository     Repository
	resumeService  resume.Service
	engineService  engine.Service
	sessionManager SessionManager
}

func NewService(
	repository Repository,
	resumeService resume.Service,
	engineService engine.Service,
	sessionManager SessionManager,
) Service {
	return &service{
		repository:     repository,
		resumeService:  resumeService,
		engineService:  engineService,
		sessionManager: sessionManager,
	}
}

// ----------------------------------------------------------------------
// Create Interview
// ----------------------------------------------------------------------

func (s *service) CreateInterview(
	ctx context.Context,
	userID string,
	request *CreateInterviewRequest,
) (*CreateInterviewResponse, error) {

	// Verify resume belongs to this user.
	_, err := s.resumeService.GetResume(ctx, request.ResumeID, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	interview := &Interview{
		ID:          uuid.NewString(),
		UserID:      userID,
		ResumeID:    request.ResumeID,
		Status:      constants.InterviewInProgress,
		StartedAt:   nil,
		CompletedAt: nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repository.Create(ctx, interview); err != nil {
		return nil, fmt.Errorf("create interview: %w", err)
	}

	return &CreateInterviewResponse{
		ID:        interview.ID,
		ResumeID:  interview.ResumeID,
		Status:    interview.Status,
		CreatedAt: interview.CreatedAt,
	}, nil
}

// ----------------------------------------------------------------------
// Start Interview
// ----------------------------------------------------------------------

func (s *service) StartInterview(
	ctx context.Context,
	interviewID string,
	userID string,
) (*engine.InterviewState, error) {

	interview, err := s.repository.GetByID(ctx, interviewID, userID)
	if err != nil {
		return nil, err
	}

	// Resume intelligence will be loaded inside the engine.
	state, err := s.engineService.InitializeInterview(
		ctx,
		interview.ID,
		interview.UserID,
		interview.ResumeID,
	)
	if err != nil {
		return nil, ErrInterviewEngineFailed
	}

	if err := s.sessionManager.CreateSession(ctx, state); err != nil {
		return nil, err
	}

	return state, nil
}

// ----------------------------------------------------------------------
// Get Interview State
// ----------------------------------------------------------------------

func (s *service) GetInterviewState(
	ctx context.Context,
	interviewID string,
) (*engine.InterviewState, error) {

	return s.sessionManager.GetSession(ctx, interviewID)
}

// ----------------------------------------------------------------------
// Update Interview State
// ----------------------------------------------------------------------

func (s *service) UpdateInterviewState(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	state.LastUpdatedAt = time.Now().UTC()

	return s.sessionManager.UpdateSession(ctx, state)
}

// ----------------------------------------------------------------------
// Complete Interview
// ----------------------------------------------------------------------

func (s *service) CompleteInterview(
	ctx context.Context,
	interviewID string,
	userID string,
) error {

	if err := s.repository.UpdateStatus(
		ctx,
		interviewID,
		userID,
		constants.InterviewCompleted,
	); err != nil {
		return err
	}

	if err := s.sessionManager.DeleteSession(ctx, interviewID); err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------------
// Get Interview
// ----------------------------------------------------------------------

func (s *service) GetInterview(
	ctx context.Context,
	interviewID string,
	userID string,
) (*InterviewResponse, error) {

	interview, err := s.repository.GetByID(ctx, interviewID, userID)
	if err != nil {
		return nil, err
	}

	return &InterviewResponse{
		ID:          interview.ID,
		ResumeID:    interview.ResumeID,
		Status:      interview.Status,
		StartedAt:   interview.StartedAt,
		CompletedAt: interview.CompletedAt,
		CreatedAt:   interview.CreatedAt,
		UpdatedAt:   interview.UpdatedAt,
	}, nil
}

// ----------------------------------------------------------------------
// List Interviews
// ----------------------------------------------------------------------

func (s *service) ListInterviews(
	ctx context.Context,
	userID string,
) (*ListInterviewsResponse, error) {

	interviews, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list interviews: %w", err)
	}

	response := &ListInterviewsResponse{
		Interviews: make([]InterviewResponse, 0, len(interviews)),
	}

	for _, interview := range interviews {
		response.Interviews = append(response.Interviews, InterviewResponse{
			ID:          interview.ID,
			ResumeID:    interview.ResumeID,
			Status:      interview.Status,
			StartedAt:   interview.StartedAt,
			CompletedAt: interview.CompletedAt,
			CreatedAt:   interview.CreatedAt,
			UpdatedAt:   interview.UpdatedAt,
		})
	}

	return response, nil
}