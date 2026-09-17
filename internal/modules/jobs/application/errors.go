package application

import "errors"

var (
	ErrJobNotFound       = errors.New("job not found")
	ErrDuplicateJob      = errors.New("job already exists")
)