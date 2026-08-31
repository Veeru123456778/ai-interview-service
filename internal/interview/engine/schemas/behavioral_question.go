package schemas

import "fmt"

// BehavioralQuestionOutput is the structured response returned by Gemini
// for generating a behavioral interview question.
type BehavioralQuestionOutput struct {
	QuestionID string `json:"question_id"`

	Question string `json:"question"`

	Competency string `json:"competency"`

	FollowUpAllowed bool `json:"follow_up_allowed"`

	EvaluationFocus []string `json:"evaluation_focus"`
}

// ValidateBehavioralQuestionOutput validates Gemini output.
func ValidateBehavioralQuestionOutput(
	output *BehavioralQuestionOutput,
) error {

	if output.QuestionID == "" {
		return fmt.Errorf("question_id is required")
	}

	if output.Question == "" {
		return fmt.Errorf("question is required")
	}

	if output.Competency == "" {
		return fmt.Errorf("competency is required")
	}

	if len(output.EvaluationFocus) == 0 {
		return fmt.Errorf("evaluation_focus cannot be empty")
	}

	return nil
}