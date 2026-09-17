package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rashad-Muntar/forge/internal/modules/jobs/domain"
)

type CreateJobCommand struct {
	IdempotencyKey string
	Name           string
	Type           string
	Payload        []byte
	Priority       int
	MaxAttempts    int
}

type JobService struct {
	repository JobRepository
}

func NewJobService(repository JobRepository) *JobService {
	return &JobService{
		repository: repository,
	}
}

func (s *JobService) CreateJob(
	ctx context.Context,
	command CreateJobCommand,
) (*domain.Job, error) {
	existingJob, err := s.repository.GetByIdempotencyKey(
		ctx,
		command.IdempotencyKey,
	)

	if err == nil {
		return existingJob, nil
	}

	if !errors.Is(err, ErrJobNotFound) {
		return nil, fmt.Errorf("check idempotency key: %w", err)
	}

	job, err := domain.NewJob(domain.CreateJobParams{
		IdempotencyKey: command.IdempotencyKey,
		Name:           command.Name,
		Type:           command.Type,
		Payload:        command.Payload,
		Priority:       command.Priority,
		MaxAttempts:    command.MaxAttempts,
	})

	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("persist job: %w", err)
	}

	return job, nil
}