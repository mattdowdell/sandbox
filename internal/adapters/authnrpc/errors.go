package authnrpc

import (
	"errors"

	"connectrpc.com/connect"
)

// Common errors to be returned by handler methods.
var (
	ErrUnauthenticated = newError(connect.CodeUnauthenticated, "missing or invalid authentication")
	ErrInternal        = newError(connect.CodeInternal, "internal error")
	ErrUnimplemented   = newError(connect.CodeUnimplemented, "unimplemented")
)

// newError creates an error response using the given code and message.
func newError(code connect.Code, msg string) *connect.Error {
	return connect.NewError(code, errors.New(msg))
}
