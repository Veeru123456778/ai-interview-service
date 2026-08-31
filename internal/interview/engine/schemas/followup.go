package schemas

import "fmt"

// FollowUpOutput is returned by the Follow-up Generator prompt.
type FollowUpOutput struct {
	Question      string `json:"question"`
	ExpectedFocus string `json:"expected_focus"`
	Reason        string `json:"reason"`
}

func ValidateFollowUpOutput(output *FollowUpOutput) error {

	if output.Question == "" {
		return fmt.Errorf("question is required")
	}

	if output.ExpectedFocus == "" {
		return fmt.Errorf("expected_focus is required")
	}

	if output.Reason == "" {
		return fmt.Errorf("reason is required")
	}

	return nil
}