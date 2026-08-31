package schemas

import "fmt"

// ThinkingOutput is returned by the Thinking Prompt.
type ThinkingOutput struct {
	Response string `json:"response"`
}

func ValidateThinkingOutput(output *ThinkingOutput) error {

	if output.Response == "" {
		return fmt.Errorf("response is required")
	}

	return nil
}