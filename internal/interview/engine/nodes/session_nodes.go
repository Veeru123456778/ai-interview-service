package nodes

import (
	"context"
	"fmt"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
)

// ----------------------------------------------------------------------
// Initialize Session Node
// ----------------------------------------------------------------------

type initializeSessionNode struct {
	repository engine.Repository
}

func NewInitializeSessionNode(
	repository engine.Repository,
) InitializeSessionNode {
	return &initializeSessionNode{
		repository: repository,
	}
}

func (n *initializeSessionNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	now := time.Now().UTC()

	state.Status = constants.InterviewInProgress
	state.CurrentQuestionNo = 0
	state.FollowUpCount = 0

	state.AskedQuestions = []string{}
	state.CompletedTopics = []string{}
	state.ConversationHistory = []engine.ConversationMessage{}

	state.StartedAt = now
	state.LastUpdatedAt = now

	return n.repository.SaveState(ctx, state)
}

// ----------------------------------------------------------------------
// Resolve Resume Context Node
// ----------------------------------------------------------------------

type resolveResumeContextNode struct{}

func NewResolveResumeContextNode() ResolveResumeContextNode {
	return &resolveResumeContextNode{}
}

func (n *resolveResumeContextNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if len(state.ResumeContexts) == 0 {
		return fmt.Errorf("resume contains no interview contexts")
	}

	context := state.ResumeContexts[0]

	state.CurrentContextID = context.ContextID
	state.CurrentContextType = context.ContextType
	state.CurrentContextTitle = context.ContextName

	if len(context.TopicIDs) == 0 {
		return fmt.Errorf("interview context contains no topics")
	}

	state.CurrentTopicID = context.TopicIDs[0]
	state.CurrentDifficulty = constants.DifficultyEasy

	return nil
}