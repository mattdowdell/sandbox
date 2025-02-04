package main

import (
	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/adapters/healthrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/reflectrpc"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/validatex"
)

// collectHandlers merges multiple Handler implementations into a slice.
//
// While wire can cast a struct to an interface, it gets confused if multiple instances of a type
// are present. For more details, see https://github.com/google/wire/issues/207.
func collectHandlers(
	example *examplerpc.Handler,
	reflect *reflectrpc.Handler,
	health *healthrpc.Handler,
) []rpcserver.Handler {
	return []rpcserver.Handler{
		example,
		reflect,
		health,
	}
}

// ...
//
// TODO: document that the ordering here is important.
func collectInterceptors(
	validate *validatex.Interceptor,
	otelconnect *otelconnectx.Interceptor,
) []connect.Interceptor {
	return []connect.Interceptor{
		otelconnect,
		validate,
	}
}

// ...
//
// TODO: document implementations of ordering with regards to panic recovery.
func collectHandlerOptions(
	interceptors []connect.Interceptor,
	recoverer *rpcserver.Recoverer,
) []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
		connect.WithRecover(recoverer.Handle),
	}
}
