package interview

import "time"

// Request to create a new interview from a resume.
type CreateInterviewRequest struct {
	ResumeID string `json:"resume_id" binding:"required,uuid"`
}

// Response returned after creating an interview.
type CreateInterviewResponse struct {
	ID        string    `json:"id"`
	ResumeID  string    `json:"resume_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Response returned when fetching interview details.
type InterviewResponse struct {
	ID          string     `json:"id"`
	ResumeID    string     `json:"resume_id"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Response returned when listing interviews for a user.
type ListInterviewsResponse struct {
	Interviews []InterviewResponse `json:"interviews"`
}