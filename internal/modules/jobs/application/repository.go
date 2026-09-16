package application

import (
	"context"

	"github.com/Rashad-Muntar/forge/internal/modules/jobs/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job *domain.Job) error
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Job, error)
}