package nodes

import (
	"context"
	"fmt"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/constants"
)

// ----------------------------------------------------------------------
// Behavioral Discussion Node
// ----------------------------------------------------------------------

type behavioralDiscussionNode struct{}

func NewBehavioralDiscussionNode() BehavioralDiscussionNode {
	return &behavioralDiscussionNode{}
}

func (n *behavioralDiscussionNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	state.CurrentContextID = "behavioral"
	state.CurrentContextType = "BEHAVIORAL"
	state.CurrentContextTitle = "Behavioral Discussion"

	state.CurrentTopicID = "behavioral"
	state.CurrentScenario = "BEHAVIORAL"

	state.CurrentQuestionID = ""
	state.FollowUpCount = 0

	return nil
}

// ----------------------------------------------------------------------
// Generate Final Evaluation Node
// ----------------------------------------------------------------------

type generateFinalEvaluationNode struct{}

func NewGenerateFinalEvaluationNode() GenerateFinalEvaluationNode {
	return &generateFinalEvaluationNode{}
}

func (n *generateFinalEvaluationNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	_ = ctx

	if len(state.TopicScores) == 0 {
		return fmt.Errorf("no topic scores available")
	}

	totalScore := 0.0
	totalQuestions := 0

	for _, topic := range state.TopicScores {
		totalScore += topic.Score
		totalQuestions += topic.QuestionsAsked
	}

	if totalQuestions == 0 {
		return fmt.Errorf("no questions were evaluated")
	}

	averageScore := totalScore / float64(len(state.TopicScores))

	recommendation := constants.RecommendationNoHire

	switch {
	case averageScore >= 8:
		recommendation = constants.RecommendationStrongHire

	case averageScore >= 6.5:
		recommendation = constants.RecommendationHire

	case averageScore >= 5:
		recommendation = constants.RecommendationHold
	}

	state.FinalEvaluation = &engine.FinalEvaluation{
		AverageScore:   averageScore,
		Recommendation: recommendation,
		TotalQuestions: totalQuestions,
	}

	return nil
}