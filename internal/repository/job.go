package repository

import (
	"context"
	"fmt"

	"mpesa-finance/internal/database"
	"mpesa-finance/internal/models"

	"github.com/jackc/pgx/v5"
)

type JobRepository struct {
	db *database.DB
}

func NewJobRepository(db *database.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (id, user_id, file_path, original_filename, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	return r.db.Pool.QueryRow(
		ctx, query,
		job.ID,
		job.UserID,
		job.FilePath,
		job.OriginalFilename,
		job.Status,
	).Scan(&job.CreatedAt, &job.UpdatedAt)
}

func (r *JobRepository) GetByID(ctx context.Context, jobID string) (*models.Job, error) {
	query := `
		SELECT id, user_id, file_path, original_filename, status, 
		       error_message, created_at, updated_at, completed_at
		FROM jobs
		WHERE id = $1
	`

	job := &models.Job{}
	err := r.db.Pool.QueryRow(ctx, query, jobID).Scan(
		&job.ID,
		&job.UserID,
		&job.FilePath,
		&job.OriginalFilename,
		&job.Status,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.CompletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("job not found")
	}
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (r *JobRepository) GetByUserID(ctx context.Context, userID string, limit int) ([]*models.Job, error) {
	query := `
		SELECT id, user_id, file_path, original_filename, status,
		       error_message, created_at, updated_at, completed_at
		FROM jobs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		err := rows.Scan(
			&job.ID,
			&job.UserID,
			&job.FilePath,
			&job.OriginalFilename,
			&job.Status,
			&job.ErrorMessage,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

func (r *JobRepository) UpdateStatus(ctx context.Context, jobID string, status models.JobStatus, errorMessage string) error {
	query := `
    UPDATE jobs
    SET status = $1::job_status, 
        error_message = $2,
        updated_at = NOW(),
        completed_at = CASE WHEN $1::text IN ('completed', 'failed') THEN NOW() ELSE completed_at END
    WHERE id = $3
`
	result, err := r.db.Pool.Exec(ctx, query, status, errorMessage, jobID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

func (r *JobRepository) GetNextQueuedJob(ctx context.Context) (*models.Job, error) {
	// Get the oldest queued job and mark it as processing
	query := `
		UPDATE jobs
		SET status = 'processing', updated_at = NOW()
		WHERE id = (
			SELECT id FROM jobs
			WHERE status = 'queued'
			ORDER BY created_at ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, user_id, file_path, original_filename, status,
		          error_message, created_at, updated_at, completed_at
	`

	job := &models.Job{}
	err := r.db.Pool.QueryRow(ctx, query).Scan(
		&job.ID,
		&job.UserID,
		&job.FilePath,
		&job.OriginalFilename,
		&job.Status,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.CompletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // No jobs available
	}
	if err != nil {
		return nil, err
	}

	return job, nil
}
