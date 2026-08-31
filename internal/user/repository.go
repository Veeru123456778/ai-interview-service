package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *User) error

	GetByID(ctx context.Context, userID string) (*User, error)

	GetBySupabaseUserID(
		ctx context.Context,
		supabaseUserID string,
	) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

// ----------------------------------------------------------------------
// Create User
// ----------------------------------------------------------------------

func (r *repository) Create(
	ctx context.Context,
	user *User,
) error {

	query := `
		INSERT INTO users (
			id,
			supabase_user_id,
			email,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.SupabaseUserID,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------
// Get User By ID
// ----------------------------------------------------------------------

func (r *repository) GetByID(
	ctx context.Context,
	userID string,
) (*User, error) {

	query := `
		SELECT
			id,
			supabase_user_id,
			email,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.SupabaseUserID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

// ----------------------------------------------------------------------
// Get User By Supabase User ID
// ----------------------------------------------------------------------

func (r *repository) GetBySupabaseUserID(
	ctx context.Context,
	supabaseUserID string,
) (*User, error) {

	query := `
		SELECT
			id,
			supabase_user_id,
			email,
			created_at,
			updated_at
		FROM users
		WHERE supabase_user_id = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, supabaseUserID).Scan(
		&user.ID,
		&user.SupabaseUserID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by supabase id: %w", err)
	}

	return &user, nil
}