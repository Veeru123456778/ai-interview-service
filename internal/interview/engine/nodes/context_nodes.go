package nodes

import (
	"context"
	"fmt"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
)

// ----------------------------------------------------------------------
// Select Context Node
// ----------------------------------------------------------------------

type selectContextNode struct{}

func NewSelectContextNode() SelectContextNode {
	return &selectContextNode{}
}

func (n *selectContextNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if len(state.ResumeContexts) == 0 {
		return fmt.Errorf("resume contains no interview contexts")
	}

	if state.CurrentContextIndex >= len(state.ResumeContexts) {
		return fmt.Errorf("no interview contexts remaining")
	}

	currentContext := state.ResumeContexts[state.CurrentContextIndex]

	state.CurrentContextID = currentContext.ContextID
	state.CurrentContextType = currentContext.ContextType
	state.CurrentContextTitle = currentContext.ContextName

	if len(currentContext.TopicIDs) == 0 {
		return fmt.Errorf("context %s contains no topics", currentContext.ContextID)
	}

	state.CurrentTopicID = currentContext.TopicIDs[0]
	state.CurrentDifficulty = constants.DifficultyEasy
	state.FollowUpCount = 0

	return nil
}

// ----------------------------------------------------------------------
// Transition Context Node
// ----------------------------------------------------------------------

type transitionContextNode struct{}

func NewTransitionContextNode() TransitionContextNode {
	return &transitionContextNode{}
}

func (n *transitionContextNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if len(state.ResumeContexts) == 0 {
		return fmt.Errorf("resume contains no interview contexts")
	}

	currentContext := state.ResumeContexts[state.CurrentContextIndex]

	// Mark current topic as completed.
	if state.CurrentTopicID != "" && !contains(state.CompletedTopics, state.CurrentTopicID) {
		state.CompletedTopics = append(state.CompletedTopics, state.CurrentTopicID)
	}

	// Try next topic within the same context.
	for _, topicID := range currentContext.TopicIDs {
		if !contains(state.CompletedTopics, topicID) {
			state.CurrentTopicID = topicID
			state.CurrentDifficulty = constants.DifficultyEasy
			state.FollowUpCount = 0
			return nil
		}
	}

	// Move to next context.
	state.CurrentContextIndex++

	if state.CurrentContextIndex >= len(state.ResumeContexts) {
		state.CurrentContextID = ""
		state.CurrentContextType = ""
		state.CurrentContextTitle = ""
		state.CurrentTopicID = ""
		return nil
	}

	nextContext := state.ResumeContexts[state.CurrentContextIndex]

	state.CurrentContextID = nextContext.ContextID
	state.CurrentContextType = nextContext.ContextType
	state.CurrentContextTitle = nextContext.ContextName

	if len(nextContext.TopicIDs) == 0 {
		return fmt.Errorf("context %s contains no topics", nextContext.ContextID)
	}

	state.CurrentTopicID = nextContext.TopicIDs[0]
	state.CurrentDifficulty = constants.DifficultyEasy
	state.FollowUpCount = 0

	return nil
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}