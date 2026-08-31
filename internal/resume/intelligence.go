package resume

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

/*
ResumeIntelligence is the canonical structure stored in PostgreSQL.

The Interview Engine loads this structure before starting an interview.
*/

type ResumeIntelligence struct {
	TechnologyGraph  []TechnologyNode  `json:"technology_graph"`
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

	technologyGraph := buildTechnologyGraph(output)
	interviewContexts := buildInterviewContexts(output, technologyGraph)

	return &ResumeIntelligence{
		TechnologyGraph:  technologyGraph,
		InterviewContexts: interviewContexts,
	}, nil
}

// ----------------------------------------------------------------------
// Technology Graph Builder
// ----------------------------------------------------------------------

func buildTechnologyGraph(output *ResumeParserOutput) []TechnologyNode {

	topicMap := make(map[string]TechnologyNode)

	addTechnology := func(name string) {
		name = strings.TrimSpace(name)
		if name == "" {
			return
		}

		key := strings.ToLower(name)

		if _, exists := topicMap[key]; exists {
			return
		}

		topicMap[key] = TechnologyNode{
			TopicID:    uuid.NewString(),
			Name:       name,
			Category:   inferTechnologyCategory(name),
			Confidence: 1.0,
		}
	}

	for _, skill := range output.Skills {
		addTechnology(skill.Name)
	}

	for _, project := range output.Projects {
		for _, technology := range project.Technologies {
			addTechnology(technology)
		}
	}

	for _, experience := range output.Experience {
		for _, technology := range experience.Technologies {
			addTechnology(technology)
		}
	}

	graph := make([]TechnologyNode, 0, len(topicMap))

	for _, node := range topicMap {
		graph = append(graph, node)
	}

	return graph
}

// ----------------------------------------------------------------------
// Interview Context Builder
// ----------------------------------------------------------------------

func buildInterviewContexts(
	output *ResumeParserOutput,
	graph []TechnologyNode,
) []InterviewContext {

	topicIDMap := make(map[string]string)

	for _, node := range graph {
		topicIDMap[strings.ToLower(node.Name)] = node.TopicID
	}

	contexts := make([]InterviewContext, 0)

	priority := 1

	for _, project := range output.Projects {

		topicIDs := make([]string, 0)

		for _, technology := range project.Technologies {
			if id, ok := topicIDMap[strings.ToLower(technology)]; ok {
				topicIDs = append(topicIDs, id)
			}
		}

		contexts = append(contexts, InterviewContext{
			ContextID:   uuid.NewString(),
			ContextType: "PROJECT",
			ContextName: project.Name,
			Priority:    priority,
			TopicIDs:    topicIDs,
		})

		priority++
	}

	for _, experience := range output.Experience {

		topicIDs := make([]string, 0)

		for _, technology := range experience.Technologies {
			if id, ok := topicIDMap[strings.ToLower(technology)]; ok {
				topicIDs = append(topicIDs, id)
			}
		}

		contexts = append(contexts, InterviewContext{
			ContextID:   uuid.NewString(),
			ContextType: "EXPERIENCE",
			ContextName: experience.Company,
			Priority:    priority,
			TopicIDs:    topicIDs,
		})

		priority++
	}

	return contexts
}

// ----------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------

func inferTechnologyCategory(name string) string {

	switch strings.ToLower(name) {
	case "postgresql", "mysql", "mongodb", "redis":
		return "Database"

	case "docker", "kubernetes":
		return "DevOps"

	case "aws", "gcp", "azure":
		return "Cloud"

	case "react", "next.js", "nextjs", "angular", "vue":
		return "Frontend"

	case "go", "golang", "java", "python", "javascript", "typescript":
		return "Language"

	default:
		return "Technology"
	}
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