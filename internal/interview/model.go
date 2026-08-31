package interview

import "time"

// Represents an interview session stored in PostgreSQL.

type Interview struct {
	ID string `json:"id"`

	UserID string `json:"user_id"`

	ResumeID string `json:"resume_id"`

	Status string `json:"status"`

	StartedAt *time.Time `json:"started_at,omitempty"`

	CompletedAt *time.Time `json:"completed_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}