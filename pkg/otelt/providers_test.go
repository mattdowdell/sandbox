package otelt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/mattdowdell/sandbox/pkg/otelt"
)

const (
	testPackage     = "example.com/package"
	testCounterName = "example_counter"
	testSpanName    = "example_span"
)

func Test_NewMeterProvider(t *testing.T) {
	// arrange

	// act
	provider, collect := otelt.NewMeterProvider()

	// assert
	assert.NotNil(t, provider)
	assert.NotNil(t, collect)
}

func Test_NewMeterProvider_Collect(t *testing.T) {
	// arrange
	provider, collect := otelt.NewMeterProvider()

	counter, err := provider.Meter(testPackage).Int64Counter(testCounterName)
	require.NoError(t, err)

	counter.Add(t.Context(), 1)

	// act
	got, err := collect(t.Context())

	// assert
	want := metricdata.ResourceMetrics{
		Resource: resource.Default(),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name: testPackage,
				},
				Metrics: []metricdata.Metrics{
					{
						Name: testCounterName,
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
	assert.NoError(t, err)
}

func Test_NewTracerProvider(t *testing.T) {
	// arrange

	// act
	provider, collect := otelt.NewTracerProvider()

	// assert
	assert.NotNil(t, provider)
	assert.NotNil(t, collect)
}

func Test_NewTracerProvider_Collect(t *testing.T) {
	// arrange
	provider, collect := otelt.NewTracerProvider()

	_, span := provider.Tracer(testPackage).Start(t.Context(), testSpanName)

	span.SetAttributes(attribute.Bool("example", true))
	span.AddEvent("event")

	span.End()

	// act
	got := collect()

	// assert
	attrs := []attribute.KeyValue{attribute.Bool("example", true)}
	events := []trace.Event{{Name: "event"}}

	if assert.Len(t, got, 1) {
		otelt.AssertSpanEqual(t, got[0], testSpanName, attrs, events)
	}
}
