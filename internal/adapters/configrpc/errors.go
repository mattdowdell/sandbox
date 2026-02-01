package configrpc

import (
	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/rpcerrors"
)

// Common errors to be returned by handler methods.
var (
	ErrValueNotFound = rpcerrors.New(connect.CodeNotFound, "value does not exist")
)
