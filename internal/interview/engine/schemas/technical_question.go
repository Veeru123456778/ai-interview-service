package schemas

import "fmt"

// TechnicalQuestionOutput is the structured response expected from Gemini
// for generating the next interview question.
type TechnicalQuestionOutput struct {
	QuestionID string `json:"question_id"`

	Question string `json:"question"`

	QuestionType string `json:"question_type"`
	Scenario     string `json:"scenario"`

	ExpectedDifficulty string `json:"expected_difficulty"`

	ExpectedTopics []string `json:"expected_topics"`

	FollowUpAllowed bool `json:"follow_up_allowed"`
	HintAllowed     bool `json:"hint_allowed"`
}

// ValidateTechnicalQuestionOutput validates Gemini output before
// it enters the interview engine.
func ValidateTechnicalQuestionOutput(
	output *TechnicalQuestionOutput,
) error {

	if output.QuestionID == "" {
		return fmt.Errorf("question_id is required")
	}

	if output.Question == "" {
		return fmt.Errorf("question is required")
	}

	if output.QuestionType == "" {
		return fmt.Errorf("question_type is required")
	}

	if output.Scenario == "" {
		return fmt.Errorf("scenario is required")
	}

	if output.ExpectedDifficulty == "" {
		return fmt.Errorf("expected_difficulty is required")
	}

	if len(output.ExpectedTopics) == 0 {
		return fmt.Errorf("expected_topics cannot be empty")
	}

	return nil
}