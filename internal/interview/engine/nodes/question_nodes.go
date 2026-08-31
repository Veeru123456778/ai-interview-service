package nodes

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
)

// ----------------------------------------------------------------------
// Select Scenario Node
// ----------------------------------------------------------------------

type selectScenarioNode struct{}

func NewSelectScenarioNode() SelectScenarioNode {
	return &selectScenarioNode{}
}

func (n *selectScenarioNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if state.CurrentTopicID == "" {
		return fmt.Errorf("current topic is not selected")
	}

	// Initial implementation:
	// Every new topic starts with IMPLEMENTATION.
	// Later this will rotate between debugging, scaling,
	// performance, trade-off, etc.
	state.CurrentScenario = "IMPLEMENTATION"

	return nil
}

// ----------------------------------------------------------------------
// Generate Question Node
// ----------------------------------------------------------------------

type generateQuestionNode struct{}

func NewGenerateQuestionNode() GenerateQuestionNode {
	return &generateQuestionNode{}
}

func (n *generateQuestionNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if state.CurrentTopicID == "" {
		return fmt.Errorf("current topic is not selected")
	}

	if state.CurrentScenario == "" {
		return fmt.Errorf("scenario is not selected")
	}

	questionID := uuid.NewString()

	state.CurrentQuestionID = questionID
	state.CurrentQuestionNo++

	state.AskedQuestions = append(state.AskedQuestions, questionID)

	// Gemini integration will populate the actual question later.
	state.ConversationHistory = append(
		state.ConversationHistory,
		engine.ConversationMessage{
			Role:       "INTERVIEWER",
			Content:    "",
			QuestionID: questionID,
			TopicID:    state.CurrentTopicID,
		},
	)

	return nil
}