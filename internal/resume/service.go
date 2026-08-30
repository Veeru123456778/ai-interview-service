package resume

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
	"github.com/google/uuid"
)

type Service interface {
	CreateResume(
		ctx context.Context,
		userID string,
		file multipart.File,
		fileName string,
	) (*UploadResumeResponse, error)

	GetResume(
		ctx context.Context,
		resumeID string,
		userID string,
	) (*ResumeResponse, error)

	ListResumes(
		ctx context.Context,
		userID string,
	) (*ListResumesResponse, error)
}

type service struct {
	repository   Repository
	extractor    Extractor
	normalizer   Normalizer
	parser       Parser
	intelligence IntelligenceBuilder
}

func NewService(
	repository Repository,
	extractor Extractor,
	normalizer Normalizer,
	parser Parser,
	builder IntelligenceBuilder,
) Service {
	return &service{
		repository:   repository,
		extractor:    extractor,
		normalizer:   normalizer,
		parser:       parser,
		intelligence: builder,
	}
}

// ----------------------------------------------------------------------
// Create Resume
// ----------------------------------------------------------------------

func (s *service) CreateResume(
	ctx context.Context,
	userID string,
	file multipart.File,
	fileName string,
) (*UploadResumeResponse, error) {

	if err := validatePDFFile(fileName); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	resume := &Resume{
		ID:                uuid.NewString(),
		UserID:            userID,
		FileName:          fileName,
		Status:            constants.ResumeProcessing,
		TechnologyGraph:   nil,
		InterviewContexts: nil,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Step 1: Save initial resume record.
	if err := s.repository.Create(ctx, resume); err != nil {
		return nil, fmt.Errorf("create resume: %w", err)
	}

	// ------------------------------------------------------------------
	// Resume Processing Pipeline
	// ------------------------------------------------------------------

	// Step 2: Extract text from PDF.
	rawText, err := s.extractor.ExtractText(file)
	if err != nil {
		_ = s.repository.UpdateStatus(
			ctx,
			resume.ID,
			userID,
			constants.ResumeFailed,
		)
		return nil, apperrors.ErrResumeExtractionFailed
	}

	// Step 3: Normalize extracted text.
	normalizedText := s.normalizer.Normalize(rawText)

	// Step 4: Parse resume using LLM.
	parserOutput, err := s.parser.Parse(ctx, normalizedText)
	if err != nil {
		_ = s.repository.UpdateStatus(
			ctx,
			resume.ID,
			userID,
			constants.ResumeFailed,
		)
		return nil, apperrors.ErrResumeParserFailed
	}

	// Step 5: Build Resume Intelligence.
	intelligence, err := s.intelligence.Build(parserOutput)
	if err != nil {
		_ = s.repository.UpdateStatus(
			ctx,
			resume.ID,
			userID,
			constants.ResumeFailed,
		)
		return nil, apperrors.ErrResumeProcessingFailed
	}

	technologyGraph, err := s.intelligence.MarshalTechnologyGraph(
		intelligence.TechnologyGraph,
	)
	if err != nil {
		return nil, apperrors.ErrResumeProcessingFailed
	}

	interviewContexts, err := s.intelligence.MarshalInterviewContexts(
		intelligence.InterviewContexts,
	)
	if err != nil {
		return nil, apperrors.ErrResumeProcessingFailed
	}

	// Step 6: Save Resume Intelligence.
	if err := s.repository.UpdateResumeIntelligence(
		ctx,
		resume.ID,
		userID,
		technologyGraph,
		interviewContexts,
		constants.ResumeReady,
	); err != nil {
		return nil, apperrors.ErrResumeProcessingFailed
	}

	return &UploadResumeResponse{
		ID:        resume.ID,
		FileName:  resume.FileName,
		Status:    constants.ResumeReady,
		CreatedAt: resume.CreatedAt,
	}, nil
}

// ----------------------------------------------------------------------
// Get Resume
// ----------------------------------------------------------------------

func (s *service) GetResume(
	ctx context.Context,
	resumeID string,
	userID string,
) (*ResumeResponse, error) {

	resume, err := s.repository.GetByID(ctx, resumeID, userID)
	if err != nil {
		return nil, apperrors.ErrResumeNotFound
	}

	return &ResumeResponse{
		ID:        resume.ID,
		FileName:  resume.FileName,
		Status:    resume.Status,
		CreatedAt: resume.CreatedAt,
		UpdatedAt: resume.UpdatedAt,
	}, nil
}

// ----------------------------------------------------------------------
// List User Resumes
// ----------------------------------------------------------------------

func (s *service) ListResumes(
	ctx context.Context,
	userID string,
) (*ListResumesResponse, error) {

	resumes, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list resumes: %w", err)
	}

	response := &ListResumesResponse{
		Resumes: make([]ResumeResponse, 0, len(resumes)),
	}

	for _, resume := range resumes {
		response.Resumes = append(response.Resumes, ResumeResponse{
			ID:        resume.ID,
			FileName:  resume.FileName,
			Status:    resume.Status,
			CreatedAt: resume.CreatedAt,
			UpdatedAt: resume.UpdatedAt,
		})
	}

	return response, nil
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

func validatePDFFile(fileName string) error {

	extension := strings.ToLower(filepath.Ext(fileName))

	if extension != ".pdf" {
		return apperrors.ErrInvalidResumeFormat
	}

	return nil
}