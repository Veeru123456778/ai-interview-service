package resume

import (
	"encoding/json"
	"fmt"
)

/*
ResumeIntelligence is the canonical structure stored in PostgreSQL.

The Interview Engine will load this structure before starting an interview.
*/
type ResumeIntelligence struct {
	TechnologyGraph  []TechnologyNode   `json:"technology_graph"`
	InterviewContexts []InterviewContext `json:"interview_contexts"`
}

/*
IntelligenceBuilder converts validated ResumeParserOutput into
ResumeIntelligence and JSONB payloads.
*/
type IntelligenceBuilder interface {
	Build(output *ResumeParserOutput) (*ResumeIntelligence, error)

	MarshalTechnologyGraph(graph []TechnologyNode) ([]byte, error)

	MarshalInterviewContexts(contexts []InterviewContext) ([]byte, error)
}

type intelligenceBuilder struct{}

func NewIntelligenceBuilder() IntelligenceBuilder {
	return &intelligenceBuilder{}
}

// ----------------------------------------------------------------------
// Build Resume Intelligence
// ----------------------------------------------------------------------

func (b *intelligenceBuilder) Build(
	output *ResumeParserOutput,
) (*ResumeIntelligence, error) {

	if err := ValidateResumeParserOutput(output); err != nil {
		return nil, err
	}

	intelligence := &ResumeIntelligence{
		TechnologyGraph:  output.TechnologyGraph,
		InterviewContexts: output.InterviewContexts,
	}

	return intelligence, nil
}

// ----------------------------------------------------------------------
// Marshal Technology Graph for PostgreSQL JSONB
// ----------------------------------------------------------------------

func (b *intelligenceBuilder) MarshalTechnologyGraph(
	graph []TechnologyNode,
) ([]byte, error) {

	data, err := json.Marshal(graph)
	if err != nil {
		return nil, fmt.Errorf("marshal technology graph: %w", err)
	}

	return data, nil
}

// ----------------------------------------------------------------------
// Marshal Interview Contexts for PostgreSQL JSONB
// ----------------------------------------------------------------------

func (b *intelligenceBuilder) MarshalInterviewContexts(
	contexts []InterviewContext,
) ([]byte, error) {

	data, err := json.Marshal(contexts)
	if err != nil {
		return nil, fmt.Errorf("marshal interview contexts: %w", err)
	}

	return data, nil
}