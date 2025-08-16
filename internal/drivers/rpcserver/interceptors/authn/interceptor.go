package authn

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	authnv1 "github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Non-allocating compile-time check for interface compliance.
var _ connect.Interceptor = (*Interceptor)(nil)

var (
	errUnavailable = connect.NewError(connect.CodeUnavailable, errors.New("service unavailable"))
	errInternal    = connect.NewError(connect.CodeInternal, errors.New("internal error"))
)

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

	client  authnv1connect.AuthnServiceClient
	ignores map[string]struct{}
}

// New creates an Interceptor.
func New(client authnv1connect.AuthnServiceClient, options ...Option) *Interceptor {
	i := &Interceptor{
		client:  client,
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

		subject, err := i.authenticate(ctx, req.Spec().Procedure, req.Header())
		if err != nil {
			return nil, err
		}

		ctx = jwtx.SubjectIntoContext(ctx, subject)
		return next(ctx, req)
	}
}

// WrapStreamingHandler implements a server streaming request interceptor.
func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		subject, err := i.authenticate(ctx, conn.Spec().Procedure, conn.RequestHeader())
		if err != nil {
			return err
		}

		ctx = jwtx.SubjectIntoContext(ctx, subject)
		return next(ctx, conn)
	}
}

// authenticate implements the common authentication logic for both unary and streaming handlers.
func (i *Interceptor) authenticate(
	ctx context.Context,
	procedure string,
	headers http.Header,
) (string, error) {
	span := trace.SpanFromContext(ctx)
	logger := slogx.FromContext(ctx)

	service, _ := rpcserver.SplitProcedure(procedure)

	if _, ok := i.ignores[service]; ok {
		span.AddEvent("authentication skipped")
		logger.DebugContext(ctx, "authentication skipped")

		return "", nil
	}

	token, err := bearerToken(headers)
	if err != nil {
		span.SetStatus(codes.Error, "invalid authorization")
		span.RecordError(err)

		logger.DebugContext(ctx, "invalid authorization", slogx.Err(err))

		return "", connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := i.client.Authenticate(ctx, connect.NewRequest(&authnv1.AuthenticateRequest{
		Token: token,
	}))
	if err != nil {
		span.RecordError(err)

		//nolint:exhaustive // not all codes are returned and they'd map to internal anyway
		switch connect.CodeOf(err) {
		case connect.CodeUnauthenticated:
			span.SetStatus(codes.Error, "failed authentication")
			logger.DebugContext(ctx, "failed authentication", slogx.Err(err))

			return "", err

		case connect.CodeUnavailable:
			span.SetStatus(codes.Error, "authentication unavailable")
			logger.ErrorContext(ctx, "authentication unavailable", slogx.Err(err))

			return "", errUnavailable

		default:
			span.SetStatus(codes.Error, "authentication error")
			logger.ErrorContext(ctx, "authentication error", slogx.Err(err))

			return "", errInternal
		}
	}

	subject := resp.Msg.GetSubject()
	span.SetAttributes(attribute.String("subject", subject))

	return subject, nil
}
