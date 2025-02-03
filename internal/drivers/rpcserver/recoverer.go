package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type Recoverer struct {
	panics metric.Int64Counter
}

// ...
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

// ...
func (r *Recoverer) Handle(ctx context.Context, spec connect.Spec, _ http.Header, recovered any) error {
	service, method := splitProcedure(spec)

	span := trace.SpanFromContext(ctx)
	span.RecordError(fmt.Errorf("panic: %v", recovered))
	span.SetStatus(codes.Error, "panicked")

	logger := slogx.LoggerFromContext(ctx)
	logger.ErrorContext(
		ctx,
		"panicked",
		slogx.Panic(recovered),
		slogx.RPCService(service),
		slogx.RPCMethod(method),
		slogx.Stacktrace(),
	)

	r.panics.Add(ctx, 1, metric.WithAttributes(
		semconv.RPCService(service),
		semconv.RPCMethod(method),
	))

	return connect.NewError(connect.CodeInternal, errors.New("internal"))
}

// ...
func splitProcedure(spec connect.Spec) (service, method string) {
	name := strings.TrimLeft(spec.Procedure, "/")

	service, method, ok := strings.Cut(name, "/")
	if !ok {
		return "", service
	}

	return service, method
}
