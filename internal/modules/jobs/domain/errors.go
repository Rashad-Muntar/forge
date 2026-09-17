package domain

import "errors"

var (
	ErrInvalidJobName        = errors.New("job name is required")
	ErrInvalidJobType        = errors.New("job type is required")
	ErrInvalidPayload        = errors.New("job payload is required")
	ErrInvalidIdempotencyKey = errors.New("idempotency key is required")
	ErrInvalidMaxAttempts    = errors.New("max attempts must be greater than zero")
	ErrInvalidPriority       = errors.New("priority cannot be negative")
)
