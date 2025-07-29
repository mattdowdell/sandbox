package authnrpc

import (
	"errors"

	"connectrpc.com/connect"
)

// ...
var (
	ErrUnauthenticated = newError(connect.CodeUnauthenticated, "missing or invalid authentication")
	ErrInternal        = newError(connect.CodeInternal, "internal error")
	ErrUnimplemented   = newError(connect.CodeUnimplemented, "unimplemented")
)

// ...
func newError(code connect.Code, msg string) *connect.Error {
	return connect.NewError(code, errors.New(msg))
}
