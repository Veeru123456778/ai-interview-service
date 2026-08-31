package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	SyncUser(
		ctx context.Context,
		supabaseUserID string,
		email string,
	) (*SyncUserResponse, error)

	GetUser(
		ctx context.Context,
		userID string,
	) (*UserResponse, error)

	GetUserBySupabaseUserID(
		ctx context.Context,
		supabaseUserID string,
	) (*User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

// ----------------------------------------------------------------------
// Sync User
// ----------------------------------------------------------------------


func (s *service) SyncUser(
	ctx context.Context,
	supabaseUserID string,
	email string,
) (*SyncUserResponse, error) {

	// User already exists.
	existingUser, err := s.repository.GetBySupabaseUserID(ctx, supabaseUserID)
	if err == nil {
		return &SyncUserResponse{
			User: UserResponse{
				ID:             existingUser.ID,
				SupabaseUserID: existingUser.SupabaseUserID,
				Email:          existingUser.Email,
				CreatedAt:      existingUser.CreatedAt,
			},
		}, nil
	}

	// Only create a new user if no record exists.
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get user by supabase id: %w", err)
	}

	now := time.Now().UTC()

	newUser := &User{
		ID:             uuid.NewString(),
		SupabaseUserID: supabaseUserID,
		Email:          email,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repository.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &SyncUserResponse{
		User: UserResponse{
			ID:             newUser.ID,
			SupabaseUserID: newUser.SupabaseUserID,
			Email:          newUser.Email,
			CreatedAt:      newUser.CreatedAt,
		},
	}, nil
}

// ----------------------------------------------------------------------
// Get User
// ----------------------------------------------------------------------

func (s *service) GetUser(
	ctx context.Context,
	userID string,
) (*UserResponse, error) {

	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &UserResponse{
		ID:             user.ID,
		SupabaseUserID: user.SupabaseUserID,
		Email:          user.Email,
		CreatedAt:      user.CreatedAt,
	}, nil
}


func (s *service) GetUserBySupabaseUserID(
    ctx context.Context,
    supabaseUserID string,
) (*User, error) {

    return s.repository.GetBySupabaseUserID(ctx, supabaseUserID)
}