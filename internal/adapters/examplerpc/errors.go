package examplerpc

import (
	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/rpcerrors"
)

// Common errors to be returned by handler methods.
var (
	ErrResourceNotFound      = rpcerrors.New(connect.CodeNotFound, "resource does not exist")
	ErrResourceAlreadyExists = rpcerrors.New(connect.CodeAlreadyExists, "resource name already in use")
)
