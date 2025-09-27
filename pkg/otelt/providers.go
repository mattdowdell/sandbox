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

// NewMeterProvider provides a [metric.MeterProvider] suitable for use in unit tests.
//
// TODO: discuss the collect function.
//
// [metric.MeterProvider]: # todo
func NewMeterProvider() (metric.MeterProvider, MetricCollectFn) {
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

// NewTracerProvider provides a [trace.TracerProvider] suitable for use in unit tests.
//
// TODO: discuss the collect function
//
// [trace.TracerProvider]: #todo
func NewTracerProvider() (trace.TracerProvider, TraceCollectFn) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	return provider, recorder.Ended
}
