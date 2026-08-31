package user

import "time"

// Represents a user stored in PostgreSQL.

type User struct {
	ID string `json:"id"`

	SupabaseUserID string `json:"supabase_user_id"`

	Email string `json:"email"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}