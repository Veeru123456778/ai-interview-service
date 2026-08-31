package provider

import "context"

// LLMProvider is the common interface used across the project.
//
// It is consumed by:
//   - Resume Parser
//   - Interview Engine
//   - Evaluation Engine
type LLMProvider interface {

	// GenerateStructuredOutput asks the model to produce JSON that can be
	// unmarshaled directly into the supplied output struct.
	GenerateStructuredOutput(
		ctx context.Context,
		promptName string,
		input map[string]any,
		output any,
	) error

	// GenerateText asks the model to generate free-form text.
	GenerateText(
		ctx context.Context,
		promptName string,
		input map[string]any,
	) (string, error)
}