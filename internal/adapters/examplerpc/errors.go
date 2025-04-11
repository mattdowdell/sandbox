package examplerpc

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

// ...
var (
	ErrInternal      = connect.NewError(connect.CodeInternal, errors.New("internal error"))
	ErrUnimplemented = connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
)

// ...
func ErrResourceNotFound(id uuid.UUID) error {
	err := fmt.Errorf("resource not found: %s", id)
	return connect.NewError(connect.CodeNotFound, err)
}

// ...
func ErrResourceAlreadyExists(name string) error {
	err := fmt.Errorf("resource name in use: %s", name)
	return connect.NewError(connect.CodeAlreadyExists, err)
}
