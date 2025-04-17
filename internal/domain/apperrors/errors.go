package apperrors

import (
	"errors"
)

// Errors to be returned by usecases and repository implementations.
//
// Repositories should wrap these errors with sufficient context to enable debugging. Usecases are
// then responsible for recording the error using logging and/or tracing. Once recorded, usecases
// should return an unwrapped error with unambiguous meaning for the calling adapter to convert to
// an appropriate user-facing error response.
//
// Additional errors should be defined here if the meaning of a particular error cannot be mapped to
// a specific error response.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInternal      = errors.New("internal error")
	ErrUnavailable   = errors.New("unavailable")
)
