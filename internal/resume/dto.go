package resume

import "time"

// UploadResumeResponse is returned after a resume upload request succeeds.
type UploadResumeResponse struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ResumeResponse is used when returning resume details.
type ResumeResponse struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListResumesResponse is used for GET /api/v1/resumes.
type ListResumesResponse struct {
	Resumes []ResumeResponse `json:"resumes"`
}