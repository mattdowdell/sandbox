package otelt_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/mattdowdell/sandbox/pkg/otelt"
)

func ExampleNewMeterProvider() {
	t := new(testing.T)

	provider, collect := otelt.NewMeterProvider()

	counter, err := provider.Meter("path/to/package").Int64Counter("my_counter")
	if err != nil {
		t.Fatal(err)
	}

	counter.Add(t.Context(), 1)

	got, err := collect(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	want := metricdata.ResourceMetrics{
		Resource: resource.Default(),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name: "path/to/package",
				},
				Metrics: []metricdata.Metrics{
					{
						Name: "my_counter",
						Data: metricdata.Sum[int64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Value: 1,
								},
							},
						},
					},
				},
			},
		},
	}

	metricdatatest.AssertEqual(t, want, got, metricdatatest.IgnoreTimestamp())
}

func ExampleNewTracerProvider() {
	t := new(testing.T)

	provider, collect := otelt.NewTracerProvider()

	_, span := provider.Tracer(testPackage).Start(t.Context(), testSpanName)

	span.SetAttributes(attribute.Bool("example", true))
	span.AddEvent("event")

	span.End()

	got := collect()
	if len(got) != 1 {
		t.Fatal("unexpected length:", len(got))
	}

	attrs := []attribute.KeyValue{attribute.Bool("example", true)}
	events := []trace.Event{{Name: "event"}}

	otelt.AssertSpanEqual(t, got[0], testSpanName, attrs, events)
}
