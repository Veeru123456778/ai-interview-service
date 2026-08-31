package schemas

import "fmt"

// BehavioralAnswerOutput is the structured evaluation returned by Gemini
// after analyzing a candidate's behavioral answer.
type BehavioralAnswerOutput struct {
	Score int `json:"score"`

	Decision string `json:"decision"`

	Competency string `json:"competency"`

	Strengths []string `json:"strengths"`
	Gaps      []string `json:"gaps"`

	Reasoning string `json:"reasoning"`

	FollowUpRequired bool   `json:"follow_up_required"`
	FollowUpReason   string `json:"follow_up_reason"`

	Confidence float64 `json:"confidence"`
}

// ValidateBehavioralAnswerOutput validates Gemini output before it is used
// by the interview engine.
func ValidateBehavioralAnswerOutput(
	output *BehavioralAnswerOutput,
) error {

	if output.Score < 0 || output.Score > 10 {
		return fmt.Errorf("score must be between 0 and 10")
	}

	switch output.Decision {
	case "PASS", "PARTIAL_PASS", "FAIL":
	default:
		return fmt.Errorf("invalid decision")
	}

	if output.Competency == "" {
		return fmt.Errorf("competency is required")
	}

	if output.Reasoning == "" {
		return fmt.Errorf("reasoning is required")
	}

	if output.Confidence < 0 || output.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1")
	}

	return nil
}