package resume

import (
	"context"
	"encoding/json"
	"fmt"
)

// LLMProvider is implemented by internal/interview/engine/provider.
// The Resume module depends only on this interface.
type LLMProvider interface {
	GenerateStructuredOutput(
		ctx context.Context,
		promptName string,
		input map[string]any,
		output any,
	) error
}

type Parser interface {
	Parse(
		ctx context.Context,
		normalizedResume string,
	) (*ResumeParserOutput, error)
}

type parser struct {
	llm LLMProvider
}

func NewParser(llm LLMProvider) Parser {
	return &parser{
		llm: llm,
	}
}

// ----------------------------------------------------------------------
// Parse Resume using Gemini
// ----------------------------------------------------------------------

func (p *parser) Parse(
	ctx context.Context,
	normalizedResume string,
) (*ResumeParserOutput, error) {

	if normalizedResume == "" {
		return nil, fmt.Errorf("normalized resume text is empty")
	}

	output := &ResumeParserOutput{}

	err := p.llm.GenerateStructuredOutput(
		ctx,
		"resume_parser", // resolved later through prompt registry.
		map[string]any{
			"resume_text": normalizedResume,
		},
		output,
	)

	if err != nil {
		return nil, fmt.Errorf("resume parser prompt failed: %w", err)
	}

	if err := ValidateResumeParserOutput(output); err != nil {
		return nil, err
	}

	return output, nil
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

// MarshalParserOutput is useful for debugging and tests.
func MarshalParserOutput(output *ResumeParserOutput) ([]byte, error) {
	return json.MarshalIndent(output, "", "  ")
}