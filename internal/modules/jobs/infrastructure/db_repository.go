package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Rashad-Muntar/forge/internal/modules/jobs/application"
	"github.com/Rashad-Muntar/forge/internal/modules/jobs/domain"
)

type PostgresJobRepository struct {
	db *pgxpool.Pool
}

func NewPostgresJobRepository(db *pgxpool.Pool) *PostgresJobRepository {
	return &PostgresJobRepository{
		db: db,
	}
}

func (r *PostgresJobRepository) Create(
	ctx context.Context,
	job *domain.Job,
) error {
	const query = `
		INSERT INTO jobs (id, idempotency_key, name, type, payload, priority, status, attempts, max_attempts, available_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8, $9, $10, $11, $12 )
	`

	_, err := r.db.Exec(
		ctx, query, job.ID, job.IdempotencyKey, job.Name, job.Type, job.Payload,
		job.Priority, job.Status, job.Attempts, job.MaxAttempts, job.AvailableAt,
		job.CreatedAt, job.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	return nil
}

func (r *PostgresJobRepository) GetByIdempotencyKey(
	ctx context.Context,
	key string,
) (*domain.Job, error) {
	const query = `
		SELECT id, idempotency_key, name, type, payload,
			priority, status, attempts, max_attempts, vailable_at,
			created_at, updated_at
		FROM jobs
		WHERE idempotency_key = $1
	`

	job := &domain.Job{}

	err := r.db.QueryRow(ctx, query, key).Scan(
		&job.ID, &job.IdempotencyKey, &job.Name, &job.Type, &job.Payload, &job.Priority, &job.Status,
		&job.Attempts, &job.MaxAttempts, &job.AvailableAt, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrJobNotFound
		}

		return nil, fmt.Errorf("get job by idempotency key: %w", err)
	}

	return job, nil
}