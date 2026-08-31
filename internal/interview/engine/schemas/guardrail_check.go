package schemas

import "fmt"

// GuardrailCheckOutput is returned by the Guardrail Detector prompt.
type GuardrailCheckOutput struct {
	IsSafe   bool   `json:"is_safe"`
	Category string `json:"category"`
}

func ValidateGuardrailCheckOutput(output *GuardrailCheckOutput) error {

	switch output.Category {
	case "NORMAL",
		"OFF_TOPIC",
		"PROMPT_INJECTION",
		"UNSUPPORTED":
	default:
		return fmt.Errorf("invalid guardrail category")
	}

	return nil
}