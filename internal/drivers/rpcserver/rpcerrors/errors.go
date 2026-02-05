// Package rpcerrors provides a small utility for creating ConnectRPC errors based on string
// messages instead of existing errors.
//
// This is appropriate for when the encountered error should not be returned to the client, such as
// with internal errors.
package rpcerrors

import (
	"errors"

	"connectrpc.com/connect"
)

// Common errors to be returned by handler methods. Additional errors can be created with the New
// function.
//
// The errors here are ordered by their equivalent gRPC status code value, i.e. unimplemented = 12.
var (
	ErrUnimplemented   = New(connect.CodeUnimplemented, "unimplemented")
	ErrInternal        = New(connect.CodeInternal, "internal error")
	ErrUnavailable     = New(connect.CodeUnavailable, "service unavailable")
	ErrUnauthenticated = New(connect.CodeUnauthenticated, "missing or invalid authentication")
)

// New creates an error response using the given code and message.
func New(code connect.Code, msg string) *connect.Error {
	return connect.NewError(code, errors.New(msg))
}
