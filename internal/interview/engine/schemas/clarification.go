package schemas

import "fmt"

// ClarificationOutput is returned by the Clarification Generator prompt.
type ClarificationOutput struct {
	Clarification string `json:"clarification"`
}

func ValidateClarificationOutput(output *ClarificationOutput) error {

	if output.Clarification == "" {
		return fmt.Errorf("clarification is required")
	}

	return nil
}