package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/resume"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
)

type ContextBuilder interface {
	BuildInitialState(
		ctx context.Context,
		interviewID string,
		userID string,
		resumeID string,
	) (*InterviewState, error)
}

type contextBuilder struct {
	resumeRepository resume.Repository
}

func NewContextBuilder(
	resumeRepository resume.Repository,
) ContextBuilder {
	return &contextBuilder{
		resumeRepository: resumeRepository,
	}
}

// ----------------------------------------------------------------------
// Build Initial Interview State
// ----------------------------------------------------------------------

func (b *contextBuilder) BuildInitialState(
	ctx context.Context,
	interviewID string,
	userID string,
	resumeID string,
) (*InterviewState, error) {

	resumeRecord, err := b.resumeRepository.GetByID(
		ctx,
		resumeID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("load resume: %w", err)
	}

	var technologyGraph []resume.TechnologyNode
	if err := json.Unmarshal(
		resumeRecord.TechnologyGraph,
		&technologyGraph,
	); err != nil {
		return nil, fmt.Errorf("decode technology graph: %w", err)
	}

	var interviewContexts []resume.InterviewContext
	if err := json.Unmarshal(
		resumeRecord.InterviewContexts,
		&interviewContexts,
	); err != nil {
		return nil, fmt.Errorf("decode interview contexts: %w", err)
	}

	now := time.Now().UTC()

	state := &InterviewState{
		InterviewID: interviewID,
		UserID:      userID,
		ResumeID:    resumeID,

		Status: constants.InterviewInProgress,

		CurrentDifficulty: constants.DifficultyEasy,

		AskedQuestions:     []string{},
		CompletedTopics:    []string{},
		ConversationHistory: []ConversationMessage{},
		TopicScores:        buildInitialTopicScores(technologyGraph),

		StartedAt:     now,
		LastUpdatedAt: now,
		ExpiresAt:     now.Add(60 * time.Minute),
	}

	if len(interviewContexts) > 0 && len(interviewContexts[0].TopicIDs) > 0 {
		state.CurrentTopicID = interviewContexts[0].TopicIDs[0]
	}

	return state, nil
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

func buildInitialTopicScores(
	technologyGraph []resume.TechnologyNode,
) []TopicScore {

	scores := make([]TopicScore, 0, len(technologyGraph))

	for _, topic := range technologyGraph {
		scores = append(scores, TopicScore{
			TopicID:         topic.TopicID,
			TopicName:       topic.Name,
			Difficulty:      constants.DifficultyEasy,
			QuestionsAsked:  0,
			QuestionsPassed: 0,
			Score:           0,
		})
	}

	return scores
}