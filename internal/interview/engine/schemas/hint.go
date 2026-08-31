package schemas

import "fmt"

// HintOutput is returned by the Hint Generator prompt.
type HintOutput struct {
	HintLevel int    `json:"hint_level"`
	Hint      string `json:"hint"`
}

func ValidateHintOutput(output *HintOutput) error {

	if output.HintLevel < 1 || output.HintLevel > 3 {
		return fmt.Errorf("hint_level must be between 1 and 3")
	}

	if output.Hint == "" {
		return fmt.Errorf("hint is required")
	}

	return nil
}