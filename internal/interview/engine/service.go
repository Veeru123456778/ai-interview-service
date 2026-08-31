// package engine

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
// )

// // ----------------------------------------------------------------------
// // Service Interface
// // ----------------------------------------------------------------------

// type Service interface {
// 	InitializeInterview(
// 		ctx context.Context,
// 		interviewID string,
// 		userID string,
// 		resumeID string,
// 		resumeContext *ResumeContext,
// 	) (*InterviewState, error)

// 	ProcessCandidateResponse(
// 		ctx context.Context,
// 		state *InterviewState,
// 		candidateAnswer string,
// 	) (*InterviewState, error)

// 	CompleteInterview(
// 		ctx context.Context,
// 		state *InterviewState,
// 	) (*InterviewState, error)
// }

// // ----------------------------------------------------------------------
// // Service Implementation
// // ----------------------------------------------------------------------

// type service struct {
// 	graph          Graph
// 	contextBuilder ContextBuilder
// 	repository     Repository // Redis InterviewState repository
// }

// func NewService(
// 	graph Graph,
// 	contextBuilder ContextBuilder,
// 	repository Repository,
// ) Service {
// 	return &service{
// 		graph:          graph,
// 		contextBuilder: contextBuilder,
// 		repository:     repository,
// 	}
// }

// // ----------------------------------------------------------------------
// // Initialize Interview
// // ----------------------------------------------------------------------

// func (s *service) InitializeInterview(
// 	ctx context.Context,
// 	interviewID string,
// 	userID string,
// 	resumeID string,
// 	resumeContext *ResumeContext,
// ) (*InterviewState, error) {

// 	now := time.Now().UTC()

// 	state := &InterviewState{
// 		InterviewID: interviewID,
// 		UserID:      userID,
// 		ResumeID:    resumeID,

// 		Status: constants.InterviewInProgress,

// 		CurrentDifficulty: constants.DifficultyMedium,

// 		AskedQuestions:      []string{},
// 		CompletedTopics:     []string{},
// 		ConversationHistory: []ConversationMessage{},
// 		TopicScores:         []TopicScore{},

// 		StartedAt:     now,
// 		LastUpdatedAt: now,
// 		ExpiresAt:     now.Add(2 * time.Hour),
// 	}

// 	// Copy Resume Intelligence into InterviewState.
// 	if err := s.contextBuilder.BuildInterviewContext(
// 		state,
// 		resumeContext,
// 	); err != nil {
// 		return nil, fmt.Errorf("build interview context: %w", err)
// 	}

// 	// Run initialization nodes (context selection + first question).
// 	if err := s.graph.Initialize(ctx, state); err != nil {
// 		return nil, fmt.Errorf("initialize interview graph: %w", err)
// 	}

// 	// Persist initial state to Redis.
// 	if err := s.repository.SaveState(ctx, state); err != nil {
// 		return nil, fmt.Errorf("save interview state: %w", err)
// 	}

// 	return state, nil
// }

// // ----------------------------------------------------------------------
// // Process Candidate Response
// // ----------------------------------------------------------------------

// func (s *service) ProcessCandidateResponse(
// 	ctx context.Context,
// 	state *InterviewState,
// 	candidateAnswer string,
// ) (*InterviewState, error) {

// 	state.ConversationHistory = append(
// 		state.ConversationHistory,
// 		ConversationMessage{
// 			Role:       "CANDIDATE",
// 			Content:    candidateAnswer,
// 			QuestionID: state.CurrentQuestionID,
// 			TopicID:    state.CurrentTopicID,
// 			CreatedAt:  time.Now().UTC(),
// 		},
// 	)

// 	state.LastUpdatedAt = time.Now().UTC()

// 	if err := s.graph.ProcessTurn(ctx, state); err != nil {
// 		return nil, fmt.Errorf("process interview turn: %w", err)
// 	}

// 	if err := s.repository.SaveState(ctx, state); err != nil {
// 		return nil, fmt.Errorf("save interview state: %w", err)
// 	}

// 	return state, nil
// }

// // ----------------------------------------------------------------------
// // Complete Interview
// // ----------------------------------------------------------------------

// func (s *service) CompleteInterview(
// 	ctx context.Context,
// 	state *InterviewState,
// ) (*InterviewState, error) {

// 	state.Status = constants.InterviewCompleted
// 	state.LastUpdatedAt = time.Now().UTC()

// 	if err := s.graph.GenerateFinalEvaluation(ctx, state); err != nil {
// 		return nil, fmt.Errorf("generate final evaluation: %w", err)
// 	}

// 	if err := s.repository.DeleteState(ctx, state.InterviewID); err != nil {
// 		return nil, fmt.Errorf("delete interview state: %w", err)
// 	}

// 	return state, nil
// }

package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
)

// ----------------------------------------------------------------------
// Service Interface
// ----------------------------------------------------------------------

type Service interface {
	InitializeInterview(
		ctx context.Context,
		interviewID string,
		userID string,
		resumeID string,
	) (*InterviewState, error)

	ProcessCandidateResponse(
		ctx context.Context,
		state *InterviewState,
		candidateAnswer string,
	) (*InterviewState, error)

	CompleteInterview(
		ctx context.Context,
		state *InterviewState,
	) (*InterviewState, error)
}

// ----------------------------------------------------------------------
// Service Implementation
// ----------------------------------------------------------------------

type service struct {
	graph          Graph
	contextBuilder ContextBuilder
	repository     Repository // Redis InterviewState repository
}

func NewService(
	graph Graph,
	contextBuilder ContextBuilder,
	repository Repository,
) Service {
	return &service{
		graph:          graph,
		contextBuilder: contextBuilder,
		repository:     repository,
	}
}

// ----------------------------------------------------------------------
// Initialize Interview
// ----------------------------------------------------------------------

func (s *service) InitializeInterview(
	ctx context.Context,
	interviewID string,
	userID string,
	resumeID string,
) (*InterviewState, error) {

	state, err := s.contextBuilder.BuildInitialState(
		ctx,
		interviewID,
		userID,
		resumeID,
	)
	if err != nil {
		return nil, fmt.Errorf("build initial interview state: %w", err)
	}

	if err := s.graph.Initialize(ctx, state); err != nil {
		return nil, fmt.Errorf("initialize interview graph: %w", err)
	}

	if err := s.repository.SaveState(ctx, state); err != nil {
		return nil, fmt.Errorf("save interview state: %w", err)
	}

	return state, nil
}

// ----------------------------------------------------------------------
// Process Candidate Response
// ----------------------------------------------------------------------

func (s *service) ProcessCandidateResponse(
	ctx context.Context,
	state *InterviewState,
	candidateAnswer string,
) (*InterviewState, error) {

	state.ConversationHistory = append(
		state.ConversationHistory,
		ConversationMessage{
			Role:       "CANDIDATE",
			Content:    candidateAnswer,
			QuestionID: state.CurrentQuestionID,
			TopicID:    state.CurrentTopicID,
			CreatedAt:  time.Now().UTC(),
		},
	)

	state.LastUpdatedAt = time.Now().UTC()

	if err := s.graph.ProcessTurn(ctx, state); err != nil {
		return nil, fmt.Errorf("process interview turn: %w", err)
	}

	if err := s.repository.SaveState(ctx, state); err != nil {
		return nil, fmt.Errorf("save interview state: %w", err)
	}

	return state, nil
}

// ----------------------------------------------------------------------
// Complete Interview
// ----------------------------------------------------------------------

func (s *service) CompleteInterview(
	ctx context.Context,
	state *InterviewState,
) (*InterviewState, error) {

	state.Status = constants.InterviewCompleted
	state.LastUpdatedAt = time.Now().UTC()

	if err := s.graph.GenerateFinalEvaluation(ctx, state); err != nil {
		return nil, fmt.Errorf("generate final evaluation: %w", err)
	}

	if err := s.repository.DeleteState(ctx, state.InterviewID); err != nil {
		return nil, fmt.Errorf("delete interview state: %w", err)
	}

	return state, nil
}