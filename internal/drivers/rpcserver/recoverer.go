package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Recoverer handles the recovery from panics and should be used with [connect.WithRecover].
//
// When a panic occurs, the following are performed:
//
// - A log record is produced containing the panic value and stacktrace.
// - If available, the current span's status is set to error and the panic value is recorded as an event.
// - The "rpc.server.panics" counter metric is incremented. This metric includes the "rpc.service" and "rpc.method" attributes.
//
// [connect.WithRecover]: https://pkg.go.dev/connectrpc.com/connect#WithRecover
type Recoverer struct {
	panics metric.Int64Counter
}

// NewRecoverer creates a Recoverer.
func NewRecoverer() (*Recoverer, error) {
	panics, err := otelx.Meter().Int64Counter(
		"rpc.server.panics",
		metric.WithDescription("Measures the number of panics per RPC."),
	)
	if err != nil {
		return nil, err
	}

	return &Recoverer{
		panics: panics,
	}, nil
}

// Handle records a panic via logging, tracing and metrics and returns a generic internal error
// response to be passed to the client.
func (r *Recoverer) Handle(ctx context.Context, spec connect.Spec, _ http.Header, recovered any) error {
	service, method := SplitProcedure(spec)

	span := trace.SpanFromContext(ctx)
	span.RecordError(fmt.Errorf("panic: %v", recovered))
	span.SetStatus(codes.Error, "panicked")

	logger := slogx.FromContext(ctx)
	logger.ErrorContext(
		ctx,
		"panicked",
		slogx.Panic(recovered),
		slogx.Stacktrace(),
	)

	r.panics.Add(ctx, 1, metric.WithAttributes(
		// TODO: add rpc.system; see https://github.com/connectrpc/connect-go/issues/816
		semconv.RPCService(service),
		semconv.RPCMethod(method),
	))

	return connect.NewError(connect.CodeInternal, errors.New("internal"))
}
