package otelt

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// ...
type (
	MetricCollectFn func(context.Context) (metricdata.ResourceMetrics, error)
	TraceCollectFn  func() []sdktrace.ReadOnlySpan
)

// Span creates a span using the give trace ID, span ID and sampled value, panicking on an invalid
// trace ID or span ID. It can be used to produce a span with a deterministic value and is intended
// for test code only.
func MustSpan(ctx context.Context, traceID, spanID string, sampled bool) context.Context {
	ctx, err := Span(ctx, traceID, spanID, sampled)
	if err != nil {
		panic(err)
	}

	return ctx
}

// Span creates a span using the give trace ID, span ID and sampled value. It can be used to produce
// a span with a deterministic value and is intended for test code only.
func Span(ctx context.Context, traceID, spanID string, sampled bool) (context.Context, error) {
	tID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return nil, err
	}

	sID, err := trace.SpanIDFromHex(spanID)
	if err != nil {
		return nil, err
	}

	var flags trace.TraceFlags
	if sampled {
		flags = trace.FlagsSampled
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tID,
		SpanID:     sID,
		TraceFlags: flags,
	})

	return trace.ContextWithRemoteSpanContext(ctx, spanCtx), nil
}

// ...
func MeterProvider() (metric.MeterProvider, MetricCollectFn) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	return provider, func(ctx context.Context) (metricdata.ResourceMetrics, error) {
		var collected metricdata.ResourceMetrics

		if err := reader.Collect(ctx, &collected); err != nil {
			return metricdata.ResourceMetrics{}, err
		}

		return collected, nil
	}
}

// ...
func TracerProvider() (trace.TracerProvider, TraceCollectFn) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	return provider, recorder.Ended
}
