package slogx_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx/otelt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

const (
	testTraceID = "0123456789abcdef0123456789abcdef"
	testSpanID  = "0123456789abcdef"
)

type testConfig struct {
	Foo string
	Bar bool
}

func Test_Attrs(t *testing.T) {
	ctx := otelt.MustSpan(t.Context(), testTraceID, testSpanID, true /*sampled*/)
	span := trace.SpanFromContext(ctx)

	conf := testConfig{
		Foo: "foo",
		Bar: true,
	}

	tests := map[string]struct {
		got  slog.Attr
		want slog.Attr
	}{
		"config": {
			got:  slogx.Config(conf),
			want: slog.Any("config", conf),
		},
		"error": {
			got:  slogx.Err(errors.New("example")),
			want: slog.String("error", "example"),
		},
		"health status": {
			got:  slogx.HealthStatus(grpc_health_v1.HealthCheckResponse_SERVING),
			want: slog.String("health_status", "SERVING"),
		},
		"panic": {
			got:  slogx.Panic("example"),
			want: slog.Any("panic", "example"),
		},
		"trace id": {
			got:  slogx.TraceID(span),
			want: slog.String("trace_id", testTraceID),
		},
		"span id": {
			got:  slogx.SpanID(span),
			want: slog.String("span_id", testSpanID),
		},
		"sampled": {
			got:  slogx.Sampled(span),
			want: slog.Bool("sampled", true),
		},
		"subject": {
			got:  slogx.Subject("example"),
			want: slog.String("subject", "example"),
		},
		"rpc system": {
			got:  slogx.RPCSystem("connect"),
			want: slog.String("rpc_system", "connect"),
		},
		"rpc service": {
			got:  slogx.RPCService("example"),
			want: slog.String("rpc_service", "example"),
		},
		"rpc method": {
			got:  slogx.RPCMethod("example"),
			want: slog.String("rpc_method", "example"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}
