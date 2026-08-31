package resume

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, resume *Resume) error

	GetByID(
		ctx context.Context,
		resumeID string,
		userID string,
	) (*Resume, error)

	ListByUserID(
		ctx context.Context,
		userID string,
	) ([]Resume, error)

	UpdateStatus(
		ctx context.Context,
		resumeID string,
		userID string,
		status string,
	) error

	UpdateResumeIntelligence(
		ctx context.Context,
		resumeID string,
		userID string,
		technologyGraph []byte,
		interviewContexts []byte,
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
// Create Resume
// ----------------------------------------------------------------------

func (r *repository) Create(ctx context.Context, resume *Resume) error {

	query := `
		INSERT INTO resumes (
			id,
			user_id,
			file_name,
			status,
			technology_graph,
			interview_contexts,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		resume.ID,
		resume.UserID,
		resume.FileName,
		resume.Status,
		resume.TechnologyGraph,
		resume.InterviewContexts,
		resume.CreatedAt,
		resume.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create resume: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------
// Get Resume By ID
// ----------------------------------------------------------------------

func (r *repository) GetByID(
	ctx context.Context,
	resumeID string,
	userID string,
) (*Resume, error) {

	query := `
		SELECT
			id,
			user_id,
			file_name,
			status,
			technology_graph,
			interview_contexts,
			created_at,
			updated_at
		FROM resumes
		WHERE id = $1
		  AND user_id = $2
	`

	var resume Resume

	err := r.db.QueryRow(
		ctx,
		query,
		resumeID,
		userID,
	).Scan(
		&resume.ID,
		&resume.UserID,
		&resume.FileName,
		&resume.Status,
		&resume.TechnologyGraph,
		&resume.InterviewContexts,
		&resume.CreatedAt,
		&resume.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrResumeNotFound
		}
		return nil, fmt.Errorf("get resume by id: %w", err)
	}

	return &resume, nil
}

// ----------------------------------------------------------------------
// List Resumes For User
// ----------------------------------------------------------------------

func (r *repository) ListByUserID(
	ctx context.Context,
	userID string,
) ([]Resume, error) {

	query := `
		SELECT
			id,
			user_id,
			file_name,
			status,
			technology_graph,
			interview_contexts,
			created_at,
			updated_at
		FROM resumes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list resumes: %w", err)
	}
	defer rows.Close()

	resumes := make([]Resume, 0)

	for rows.Next() {

		var resume Resume

		err := rows.Scan(
			&resume.ID,
			&resume.UserID,
			&resume.FileName,
			&resume.Status,
			&resume.TechnologyGraph,
			&resume.InterviewContexts,
			&resume.CreatedAt,
			&resume.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan resume: %w", err)
		}

		resumes = append(resumes, resume)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resumes: %w", err)
	}

	return resumes, nil
}

// ----------------------------------------------------------------------
// Update Resume Status
// ----------------------------------------------------------------------

func (r *repository) UpdateStatus(
	ctx context.Context,
	resumeID string,
	userID string,
	status string,
) error {

	query := `
		UPDATE resumes
		SET
			status = $3,
			updated_at = NOW()
		WHERE id = $1
		  AND user_id = $2
	`

	commandTag, err := r.db.Exec(
	ctx,
	query,
	resumeID,
	userID,
	status,
	)

	if err != nil {
		return fmt.Errorf("update resume status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrResumeNotFound
	}

	return nil
}

// ----------------------------------------------------------------------
// Save Resume Intelligence
// ----------------------------------------------------------------------

func (r *repository) UpdateResumeIntelligence(
	ctx context.Context,
	resumeID string,
	userID string,
	technologyGraph []byte,
	interviewContexts []byte,
	status string,
) error {

	query := `
		UPDATE resumes
		SET
			technology_graph = $3,
			interview_contexts = $4,
			status = $5,
			updated_at = NOW()
		WHERE id = $1
		  AND user_id = $2
	`

	commandTag, err := r.db.Exec(
	ctx,
	query,
	resumeID,
	userID,
	technologyGraph,
	interviewContexts,
	status,
	)

	if err != nil {
		return fmt.Errorf("update resume intelligence: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrResumeNotFound
	}

	return nil
}