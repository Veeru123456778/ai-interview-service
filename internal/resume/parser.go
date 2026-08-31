package resume

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Veeru123456778/ai-interview-service/internal/ai/provider"
)

type Parser interface {
	Parse(
		ctx context.Context,
		normalizedResume string,
	) (*ResumeParserOutput, error)
}

type parser struct {
	llm provider.LLMProvider
}

func NewParser(llm provider.LLMProvider) Parser {
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
		"internal/resume/prompts/resume_parser_v1.txt",
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