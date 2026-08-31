package schemas

import "fmt"

// FinalEvaluationOutput is the structured interview report returned by Gemini
// after the complete interview is finished.
type FinalEvaluationOutput struct {
	OverallScore float64 `json:"overall_score"`

	Recommendation string `json:"recommendation"`

	Summary string `json:"summary"`

	Strengths []string `json:"strengths"`
	Weaknesses []string `json:"weaknesses"`

	TopicEvaluations []TopicEvaluation `json:"topic_evaluations"`

	CommunicationScore float64 `json:"communication_score"`
	ProblemSolvingScore float64 `json:"problem_solving_score"`
	SystemDesignScore float64 `json:"system_design_score"`

	HiringConfidence float64 `json:"hiring_confidence"`
}

type TopicEvaluation struct {
	TopicID string `json:"topic_id"`

	TopicName string `json:"topic_name"`

	Score float64 `json:"score"`

	Feedback string `json:"feedback"`
}

// ValidateFinalEvaluationOutput validates Gemini output before it is stored
// or returned to the client.
func ValidateFinalEvaluationOutput(
	output *FinalEvaluationOutput,
) error {

	if output.OverallScore < 0 || output.OverallScore > 10 {
		return fmt.Errorf("overall_score must be between 0 and 10")
	}

	switch output.Recommendation {
	case "STRONG_HIRE", "HIRE", "HOLD", "NO_HIRE":
	default:
		return fmt.Errorf("invalid recommendation")
	}

	if output.Summary == "" {
		return fmt.Errorf("summary is required")
	}

	if output.CommunicationScore < 0 || output.CommunicationScore > 10 {
		return fmt.Errorf("communication_score must be between 0 and 10")
	}

	if output.ProblemSolvingScore < 0 || output.ProblemSolvingScore > 10 {
		return fmt.Errorf("problem_solving_score must be between 0 and 10")
	}

	if output.SystemDesignScore < 0 || output.SystemDesignScore > 10 {
		return fmt.Errorf("system_design_score must be between 0 and 10")
	}

	if output.HiringConfidence < 0 || output.HiringConfidence > 1 {
		return fmt.Errorf("hiring_confidence must be between 0 and 1")
	}

	for _, topic := range output.TopicEvaluations {
		if topic.TopicID == "" {
			return fmt.Errorf("topic_id is required")
		}

		if topic.TopicName == "" {
			return fmt.Errorf("topic_name is required")
		}

		if topic.Score < 0 || topic.Score > 10 {
			return fmt.Errorf("topic score must be between 0 and 10")
		}
	}

	return nil
}