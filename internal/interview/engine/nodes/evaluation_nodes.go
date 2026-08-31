package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
)

// ----------------------------------------------------------------------
// Analyze Technical Answer Node
// ----------------------------------------------------------------------

type analyzeTechnicalAnswerNode struct{}

func NewAnalyzeTechnicalAnswerNode() AnalyzeTechnicalAnswerNode {
	return &analyzeTechnicalAnswerNode{}
}

func (n *analyzeTechnicalAnswerNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if len(state.ConversationHistory) == 0 {
		return fmt.Errorf("conversation history is empty")
	}

	lastMessage := state.ConversationHistory[len(state.ConversationHistory)-1]

	if lastMessage.Role != "CANDIDATE" {
		return fmt.Errorf("last message is not from candidate")
	}

	answer := strings.TrimSpace(lastMessage.Content)

	if answer == "" {
		return fmt.Errorf("candidate answer is empty")
	}

	score := 0.0

	switch {
	case len(answer) > 250:
		score = 9
	case len(answer) > 150:
		score = 7
	case len(answer) > 60:
		score = 5
	default:
		score = 3
	}

	for index := range state.TopicScores {
		if state.TopicScores[index].TopicID == state.CurrentTopicID {
			state.TopicScores[index].QuestionsAsked++
			state.TopicScores[index].Score = score

			if score >= 6 {
				state.TopicScores[index].QuestionsPassed++
			}
			break
		}
	}

	return nil
}

// ----------------------------------------------------------------------
// Update Candidate Memory Node
// ----------------------------------------------------------------------

type updateCandidateMemoryNode struct{}

func NewUpdateCandidateMemoryNode() UpdateCandidateMemoryNode {
	return &updateCandidateMemoryNode{}
}

func (n *updateCandidateMemoryNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	state.LastUpdatedAt = state.ConversationHistory[len(state.ConversationHistory)-1].CreatedAt

	return nil
}

// ----------------------------------------------------------------------
// Decide Difficulty Node
// ----------------------------------------------------------------------

type decideDifficultyNode struct{}

func NewDecideDifficultyNode() DecideDifficultyNode {
	return &decideDifficultyNode{}
}

func (n *decideDifficultyNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	for _, topic := range state.TopicScores {

		if topic.TopicID != state.CurrentTopicID {
			continue
		}

		switch {
		case topic.Score >= 8:
			state.CurrentDifficulty = constants.DifficultyHard

		case topic.Score >= 5:
			state.CurrentDifficulty = constants.DifficultyMedium

		default:
			state.CurrentDifficulty = constants.DifficultyEasy
		}

		return nil
	}

	return fmt.Errorf("topic score not found for current topic")
}