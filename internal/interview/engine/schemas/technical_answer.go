package schemas

import "fmt"

// TechnicalAnswerOutput is the structured evaluation returned by Gemini
// after analyzing a candidate's answer.
type TechnicalAnswerOutput struct {
	Score int `json:"score"`

	Decision string `json:"decision"`

	Strengths []string `json:"strengths"`
	Gaps      []string `json:"gaps"`

	Reasoning string `json:"reasoning"`

	FollowUpRequired bool `json:"follow_up_required"`
	FollowUpReason   string `json:"follow_up_reason"`

	RecommendedDifficulty string `json:"recommended_difficulty"`

	Confidence float64 `json:"confidence"`
}

// ValidateTechnicalAnswerOutput validates Gemini output before
// it updates the interview state.
func ValidateTechnicalAnswerOutput(
	output *TechnicalAnswerOutput,
) error {

	if output.Score < 0 || output.Score > 10 {
		return fmt.Errorf("score must be between 0 and 10")
	}

	switch output.Decision {
	case "PASS", "PARTIAL_PASS", "FAIL":
	default:
		return fmt.Errorf("invalid decision")
	}

	switch output.RecommendedDifficulty {
	case "EASY", "MEDIUM", "HARD":
	default:
		return fmt.Errorf("invalid recommended difficulty")
	}

	if output.Confidence < 0 || output.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1")
	}

	if output.Reasoning == "" {
		return fmt.Errorf("reasoning is required")
	}

	return nil
}