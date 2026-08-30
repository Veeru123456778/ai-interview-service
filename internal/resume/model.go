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