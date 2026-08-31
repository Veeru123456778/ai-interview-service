package user

import "time"

// Response returned for the authenticated user.
type UserResponse struct {
	ID             string    `json:"id"`
	SupabaseUserID string    `json:"supabase_user_id"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
}

// Response returned after syncing/creating the authenticated user.
type SyncUserResponse struct {
	User UserResponse `json:"user"`
}