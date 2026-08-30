package resume

import (
	"fmt"

	sharedvalidator "github.com/Veeru123456778/ai-interview-service/internal/shared/validator"
)

/*
ResumeParserOutput is the exact JSON schema expected from Gemini.
This is NOT the database model.
*/
type ResumeParserOutput struct {
	CandidateName string `json:"candidate_name" validate:"required"`

	Summary string `json:"summary"`

	Skills []string `json:"skills" validate:"required,min=1,dive,required"`

	TechnologyGraph []TechnologyNode `json:"technology_graph" validate:"required,min=1,dive"`

	InterviewContexts []InterviewContext `json:"interview_contexts" validate:"required,min=1,dive"`
}

/*
TechnologyNode represents one technology extracted from the resume.
*/
type TechnologyNode struct {
	ID string `json:"id" validate:"required"`

	Name string `json:"name" validate:"required"`

	Category string `json:"category" validate:"required"`
}

/*
InterviewContext represents one interview context such as a project,
internship, or work experience.
*/
type InterviewContext struct {
	ID string `json:"id" validate:"required"`

	Name string `json:"name" validate:"required"`

	Type string `json:"type" validate:"required,oneof=PROJECT EXPERIENCE INTERNSHIP"`

	Description string `json:"description"`

	Topics []ContextTopic `json:"topics" validate:"required,min=1,dive"`
}

/*
ContextTopic links a technology to a specific interview context.
*/
type ContextTopic struct {
	TopicID string `json:"topic_id" validate:"required"`

	Scenario string `json:"scenario" validate:"required"`
}

// ----------------------------------------------------------------------
// Validation
// ----------------------------------------------------------------------

func ValidateResumeParserOutput(output *ResumeParserOutput) error {

	validate := sharedvalidator.New()

	if err := validate.Struct(output); err != nil {
		return fmt.Errorf("resume parser validation failed: %w", err)
	}

	return nil
}