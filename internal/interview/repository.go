package interview

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, interview *Interview) error

	GetByID(
		ctx context.Context,
		interviewID string,
		userID string,
	) (*Interview, error)

	ListByUserID(
		ctx context.Context,
		userID string,
	) ([]Interview, error)

	UpdateStatus(
		ctx context.Context,
		interviewID string,
		userID string,
		status string,
	) error
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
// Create Interview
// ----------------------------------------------------------------------

func (r *repository) Create(
	ctx context.Context,
	interview *Interview,
) error {

	query := `
		INSERT INTO interviews (
			id,
			user_id,
			resume_id,
			status,
			started_at,
			completed_at,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		interview.ID,
		interview.UserID,
		interview.ResumeID,
		interview.Status,
		interview.StartedAt,
		interview.CompletedAt,
		interview.CreatedAt,
		interview.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create interview: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------
// Get Interview By ID
// ----------------------------------------------------------------------

func (r *repository) GetByID(
	ctx context.Context,
	interviewID string,
	userID string,
) (*Interview, error) {

	query := `
		SELECT
			id,
			user_id,
			resume_id,
			status,
			started_at,
			completed_at,
			created_at,
			updated_at
		FROM interviews
		WHERE id = $1
		AND user_id = $2
	`

	var interview Interview

	err := r.db.QueryRow(ctx, query, interviewID, userID).Scan(
		&interview.ID,
		&interview.UserID,
		&interview.ResumeID,
		&interview.Status,
		&interview.StartedAt,
		&interview.CompletedAt,
		&interview.CreatedAt,
		&interview.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get interview by id: %w", err)
	}

	return &interview, nil
}

// ----------------------------------------------------------------------
// List Interviews For User
// ----------------------------------------------------------------------

func (r *repository) ListByUserID(
	ctx context.Context,
	userID string,
) ([]Interview, error) {

	query := `
		SELECT
			id,
			user_id,
			resume_id,
			status,
			started_at,
			completed_at,
			created_at,
			updated_at
		FROM interviews
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list interviews: %w", err)
	}
	defer rows.Close()

	interviews := make([]Interview, 0)

	for rows.Next() {

		var interview Interview

		err := rows.Scan(
			&interview.ID,
			&interview.UserID,
			&interview.ResumeID,
			&interview.Status,
			&interview.StartedAt,
			&interview.CompletedAt,
			&interview.CreatedAt,
			&interview.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan interview: %w", err)
		}

		interviews = append(interviews, interview)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate interviews: %w", err)
	}

	return interviews, nil
}

// ----------------------------------------------------------------------
// Update Interview Status
// ----------------------------------------------------------------------

func (r *repository) UpdateStatus(
	ctx context.Context,
	interviewID string,
	userID string,
	status string,
) error {

	query := `
		UPDATE interviews
		SET
			status = $3,
			updated_at = NOW()
		WHERE id = $1
		AND user_id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		interviewID,
		userID,
		status,
	)

	if err != nil {
		return fmt.Errorf("update interview status: %w", err)
	}

	return nil
}