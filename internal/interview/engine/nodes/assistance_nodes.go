package nodes

import (
	"context"
	"fmt"
	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
)

// ----------------------------------------------------------------------
// Generate Clarification Node
// ----------------------------------------------------------------------

type generateClarificationNode struct{}

func NewGenerateClarificationNode() GenerateClarificationNode {
	return &generateClarificationNode{}
}

func (n *generateClarificationNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if state.CurrentQuestionID == "" {
		return fmt.Errorf("no active interview question")
	}

	clarification := "Sure. Let me rephrase the question. Focus on explaining your approach, key design decisions, and why you chose them. You do not need to write production-ready code unless requested."

	state.ConversationHistory = append(
		state.ConversationHistory,
		engine.ConversationMessage{
			Role:       "INTERVIEWER",
			Content:    clarification,
			QuestionID: state.CurrentQuestionID,
			TopicID:    state.CurrentTopicID,
		},
	)

	return nil
}

// ----------------------------------------------------------------------
// Generate Hint Node
// ----------------------------------------------------------------------

type generateHintNode struct{}

func NewGenerateHintNode() GenerateHintNode {
	return &generateHintNode{}
}

func (n *generateHintNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if state.CurrentQuestionID == "" {
		return fmt.Errorf("no active interview question")
	}

	var hint string

	switch state.FollowUpCount {

	case 0:
		hint = "Hint 1: Start by breaking the problem into smaller components before thinking about implementation."

	case 1:
		hint = "Hint 2: Think about the main data flow and the core data structures involved."

	default:
		hint = "Hint 3: Focus on the trade-offs behind your design instead of the final solution."
	}

	state.FollowUpCount++

	state.ConversationHistory = append(
		state.ConversationHistory,
		engine.ConversationMessage{
			Role:       "INTERVIEWER",
			Content:    hint,
			QuestionID: state.CurrentQuestionID,
			TopicID:    state.CurrentTopicID,
		},
	)

	return nil
}