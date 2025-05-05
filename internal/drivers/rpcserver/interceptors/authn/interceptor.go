package authn

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Non-allocating compile-time check for interface compliance.
var _ connect.Interceptor = (*Interceptor)(nil)

// Interceptor validates that RPC requests have valid authorization. Authorization can be skipped
// for specific services when creating the interceptor.
//
// This implementation takes a different approach to [connectrpc.com/authn] which provides HTTP
// middleware to wrap the HTTP request handler. By using HTTP middleware, authorization is checked
// before the request body is read and parsed. This saves resources, and protects the server from
// malicious requests. However, it sacrifices observability added by interceptors that only run once
// the request body has been parsed. While the risks of a malicious request are significant for
// public endpoints, they may be acceptable for internal endpoints. Additionally, if observability
// is added external, such as by Istio's sidecar, using HTTP middleware may not result in a loss of
// observability.
//
// See [connectrpc/otelconnect-go#164] for whether otelconnect could provide a middleware and so
// provide observability for the authn module's middleware too.
//
// For the purposes of this service, which does not use Istio or an equivalent, and is not deployed
// anywhere, publicly or otherwise, the risk of a malicious request body is acceptable.
//
// [connectrpc/otelconnect-go#164]: https://github.com/connectrpc/otelconnect-go/issues/164
type Interceptor struct {
	connect.Interceptor

	ignores map[string]struct{}
}

// New creates an Interceptor.
func New(options ...Option) *Interceptor {
	i := &Interceptor{
		ignores: map[string]struct{}{},
	}

	for _, o := range options {
		o.apply(i)
	}

	return i
}

// WrapUnary implements a client and server unary request interceptor.
func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().IsClient {
			return next(ctx, req)
		}

		if err := i.authenticate(ctx, req.Spec().Procedure, req.Header()); err != nil {
			return nil, err
		}

		return next(ctx, req)
	}
}

// WrapStreamingHandler implements a server streaming request interceptor.
func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if err := i.authenticate(ctx, conn.Spec().Procedure, conn.RequestHeader()); err != nil {
			return err
		}

		return next(ctx, conn)
	}
}

// authenticate implements the common authentication logic for both unary and streaming handlers.
func (i *Interceptor) authenticate(
	ctx context.Context,
	procedure string,
	headers http.Header,
) error {
	span := trace.SpanFromContext(ctx)
	logger := slogx.FromContext(ctx)

	service, _ := rpcserver.SplitProcedure(procedure)

	if _, ok := i.ignores[service]; ok {
		span.AddEvent("authentication skipped")
		logger.DebugContext(ctx, "authentication skipped")

		return nil
	}

	token, err := bearerToken(headers)
	if err != nil {
		span.SetStatus(codes.Error, "invalid authorization")
		span.RecordError(err)

		logger.DebugContext(ctx, "invalid authorization", slogx.Err(err))

		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	// TODO: parse token into claims + validate
	if err := parseToken(token); err != nil {
		span.SetStatus(codes.Error, "failed to parse token")
		span.RecordError(err)

		logger.DebugContext(ctx, "failed to parse token", slogx.Err(err))

		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	return nil
}
