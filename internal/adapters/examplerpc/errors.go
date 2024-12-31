package examplerpc

import (
	"errors"

	"connectrpc.com/connect"
)

// ...
var (
	ErrInternal      = connect.NewError(connect.CodeInternal, errors.New("internal"))
	ErrUnimplemented = connect.NewError(connect.CodeUnimplemented, errors.New("unimplemented"))
)
