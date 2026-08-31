package resume

import (
	"fmt"
	"strings"
)

// ----------------------------------------------------------------------
// Resume Parser Output Schema
// ----------------------------------------------------------------------

type ResumeParserOutput struct {
	Name           *CandidateName     `json:"name,omitempty"`
	Skills         []Skill            `json:"skills"`
	Projects       []Project          `json:"projects"`
	Experience     []Experience       `json:"experience,omitempty"`
	Education      []Education        `json:"education,omitempty"`
	Certifications []Certification    `json:"certifications,omitempty"`
}

// ----------------------------------------------------------------------
// Name (Optional)
// ----------------------------------------------------------------------

type CandidateName struct {
	FullName string `json:"full_name"`
}

// ----------------------------------------------------------------------
// Skills (Required)
// ----------------------------------------------------------------------

type Skill struct {
	Name string `json:"name"`
}

// ----------------------------------------------------------------------
// Projects (Required)
// ----------------------------------------------------------------------

type Project struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	Highlights   []string `json:"highlights,omitempty"`
}

// ----------------------------------------------------------------------
// Experience (Optional)
// ----------------------------------------------------------------------

type Experience struct {
	Company        string   `json:"company"`
	Role           string   `json:"role"`
	EmploymentType string   `json:"employment_type"`
	Duration       string   `json:"duration"`
	Description    string   `json:"description"`
	Technologies   []string `json:"technologies"`
	Highlights     []string `json:"highlights,omitempty"`
}

// ----------------------------------------------------------------------
// Education (Optional)
// ----------------------------------------------------------------------

type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field"`
	Duration    string `json:"duration"`
}

// ----------------------------------------------------------------------
// Certifications (Optional)
// ----------------------------------------------------------------------

type Certification struct {
	Name         string `json:"name"`
	Issuer       string `json:"issuer"`
	IssuedAt     string `json:"issued_at,omitempty"`
	CredentialID string `json:"credential_id,omitempty"`
}

// ----------------------------------------------------------------------
// Validation
// ----------------------------------------------------------------------

func ValidateResumeParserOutput(output *ResumeParserOutput) error {

	if output == nil {
		return fmt.Errorf("resume parser output is nil")
	}

	// Skills are required.
	if len(output.Skills) == 0 {
		return fmt.Errorf("resume must contain at least one skill")
	}

	for index, skill := range output.Skills {
		if strings.TrimSpace(skill.Name) == "" {
			return fmt.Errorf("skill %d has empty name", index)
		}
	}

	// Projects are required.
	if len(output.Projects) == 0 {
		return fmt.Errorf("resume must contain at least one project")
	}

	for index, project := range output.Projects {

		if strings.TrimSpace(project.Name) == "" {
			return fmt.Errorf("project %d has empty name", index)
		}

		if len(project.Technologies) == 0 {
			return fmt.Errorf("project %s must contain at least one technology", project.Name)
		}
	}

	// Experience is optional, but validate entries if present.
	for index, experience := range output.Experience {

		if strings.TrimSpace(experience.Company) == "" {
			return fmt.Errorf("experience %d has empty company", index)
		}

		if strings.TrimSpace(experience.Role) == "" {
			return fmt.Errorf("experience %d has empty role", index)
		}
	}

	// Education is optional.
	// Certifications are optional.

	return nil
}