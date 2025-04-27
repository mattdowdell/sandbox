package examplerpc

import (
	"errors"

	"connectrpc.com/connect"
)

// ...
var (
	ErrResourceNotFound      = newError(connect.CodeNotFound, "resource does not exist")
	ErrResourceAlreadyExists = newError(connect.CodeAlreadyExists, "resource name already in use")
	ErrInternal              = newError(connect.CodeInternal, "internal error")
	ErrUnimplemented         = newError(connect.CodeUnimplemented, "unimplemented")
)

// ...
func newError(code connect.Code, msg string) *connect.Error {
	return connect.NewError(code, errors.New(msg))
}
