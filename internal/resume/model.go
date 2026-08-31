package resume

import "time"

type Resume struct {
	ID string `json:"id"`

	UserID string `json:"user_id"`

	FileName string `json:"file_name"`

	Status string `json:"status"`

	TechnologyGraph  []byte `json:"-"`
	InterviewContexts []byte `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ----------------------------------------------------------------------
// Resume Intelligence Models
// ----------------------------------------------------------------------

type TechnologyNode struct {
	TopicID    string  `json:"topic_id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
}

type InterviewContext struct {
	ContextID   string   `json:"context_id"`
	ContextType string   `json:"context_type"`
	ContextName string   `json:"context_name"`
	Priority    int      `json:"priority"`
	TopicIDs    []string `json:"topic_ids"`
}