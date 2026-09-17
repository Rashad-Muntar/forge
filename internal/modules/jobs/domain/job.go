package domain

import (
	"crypto/rand"
	"time"
)

type UUID [16]byte

func newUUID() UUID {
	var id UUID
	if _, err := rand.Read(id[:]); err != nil {
		panic(err)
	}
	id[6] = (id[6] & 0x0f) | 0x40
	id[8] = (id[8] & 0x3f) | 0x80
	return id
}

type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusSucceeded JobStatus = "succeeded"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)


type Job struct {
	ID             UUID
	IdempotencyKey string
	Name           string
	Type           string
	Payload        []byte
	Priority       int
	Status         JobStatus
	Attempts       int
	MaxAttempts    int
	AvailableAt    time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateJobParams struct {
	IdempotencyKey string
	Name           string
	Type           string
	Payload        []byte
	Priority       int
	MaxAttempts    int
}

func NewJob(params CreateJobParams) (*Job, error) {
	if params.IdempotencyKey == "" {
		return nil, ErrInvalidIdempotencyKey
	}

	if params.Name == "" {
		return nil, ErrInvalidJobName
	}

	if params.Type == "" {
		return nil, ErrInvalidJobType
	}

	if len(params.Payload) == 0 {
		return nil, ErrInvalidPayload
	}

	if params.Priority < 0 {
		return nil, ErrInvalidPriority
	}

	if params.MaxAttempts <= 0 {
		return nil, ErrInvalidMaxAttempts
	}

	now := time.Now().UTC()

	return &Job{
		ID:             newUUID(),
		IdempotencyKey: params.IdempotencyKey,
		Name:           params.Name,
		Type:           params.Type,
		Payload:        params.Payload,
		Priority:       params.Priority,
		Status:         JobStatusQueued,
		Attempts:       0,
		MaxAttempts:    params.MaxAttempts,
		AvailableAt:    now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}