package logging

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Non-allocating compile-time check for interface compliance.
var _ connect.Interceptor = (*Interceptor)(nil)

// Interceptor adds a logger to the request context for use by server handlers and other
// interceptors. The logger contains fields identifying the RPC system, service and method.
type Interceptor struct {
	connect.Interceptor
}

// New creates an Interceptor.
func New() *Interceptor {
	return &Interceptor{}
}

// WrapUnary implements a client and server unary request interceptor.
func (*Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().IsClient {
			return next(ctx, req)
		}

		system := rpcserver.ProtocolToSystem(req.Peer())
		service, method := rpcserver.SplitProcedure(req.Spec())

		logger := slogx.FromContext(ctx).With(
			slogx.RPCSystem(system),
			slogx.RPCService(service),
			slogx.RPCMethod(method),
		)

		return next(slogx.AddToContext(ctx, logger), req)
	}
}

// WrapStreamingHandler implements a server streaming request interceptor.
func (*Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		system := rpcserver.ProtocolToSystem(conn.Peer())
		service, method := rpcserver.SplitProcedure(conn.Spec())

		logger := slogx.FromContext(ctx).With(
			slogx.RPCSystem(system),
			slogx.RPCService(service),
			slogx.RPCMethod(method),
		)

		return next(slogx.AddToContext(ctx, logger), conn)
	}
}
