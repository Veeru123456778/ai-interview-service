package interview

import (
	"context"
	"fmt"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/resume"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
	"github.com/Veeru123456778/ai-interview-service/internal/user"
	"github.com/google/uuid"
)

type Service interface {
	CreateInterview(
		ctx context.Context,
		supabaseUserID string,
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
}

type service struct {
	repository    Repository
	resumeService resume.Service
	userService   user.Service
}

func NewService(
	repository Repository,
	resumeService resume.Service,
	userService user.Service,
) Service {
	return &service{
		repository:    repository,
		resumeService: resumeService,
		userService:   userService,
	}
}

// ----------------------------------------------------------------------
// Create Interview
// ----------------------------------------------------------------------

func (s *service) CreateInterview(
	ctx context.Context,
	supabaseUserID string,
	request *CreateInterviewRequest,
) (*CreateInterviewResponse, error) {

	// Resolve internal user.
	userRecord, err := s.userService.GetUserBySupabaseUserID(
		ctx,
		supabaseUserID,
	)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	// Validate that the resume belongs to this user.
	_, err = s.resumeService.GetResume(
		ctx,
		request.ResumeID,
		userRecord.ID,
	)
	if err != nil {
		return nil, apperrors.ErrResumeNotFound
	}

	now := time.Now().UTC()

	interview := &Interview{
		ID:          uuid.NewString(),
		UserID:      userRecord.ID,
		ResumeID:    request.ResumeID,
		Status:      constants.InterviewCreated,
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
// Get Interview
// ----------------------------------------------------------------------

func (s *service) GetInterview(
	ctx context.Context,
	interviewID string,
	userID string,
) (*InterviewResponse, error) {

	interview, err := s.repository.GetByID(ctx, interviewID, userID)
	if err != nil {
		return nil, apperrors.ErrInterviewNotFound
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