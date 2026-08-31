package schemas

import "fmt"

// CandidateIntentOutput is returned by the Intent Detector prompt.
type CandidateIntentOutput struct {
	Intent     string  `json:"intent"`
	Confidence float64 `json:"confidence"`
}

func ValidateCandidateIntentOutput(output *CandidateIntentOutput) error {

	switch output.Intent {
	case "ANSWER",
		"REQUEST_CLARIFICATION",
		"ASK_HINT",
		"THINKING_OUT_LOUD",
		"OFF_TOPIC",
		"PROMPT_INJECTION",
		"END_REQUEST":
	default:
		return fmt.Errorf("invalid intent")
	}

	if output.Confidence < 0 || output.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1")
	}

	return nil
}