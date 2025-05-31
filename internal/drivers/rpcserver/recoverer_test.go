package rpcserver_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
)

func Test_NewRecoverer(t *testing.T) {
	// arrange

	// act
	recoverer, err := rpcserver.NewRecoverer()

	// assert
	assert.NotNil(t, recoverer)
	assert.NoError(t, err)
}

func Test_Recoverer_Handle(t *testing.T) {
	// arrange
	provider, collect := otelx.TestMeterProvider()
	_ = collect

	recoverer, err := rpcserver.NewRecoverer(
		rpcserver.WithMeterProvider(provider),
	)
	require.NoError(t, err)

	spec := connect.Spec{
		Procedure: "acme.foo.v1.FooService/Bar",
	}

	// act
	err = recoverer.Handle(t.Context(), spec, http.Header{}, "example")

	// assert
	require.EqualError(t, err, "internal: internal error")

	want := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:    "github.com/mattdowdell/sandbox/internal/drivers/rpcserver",
			Version: "(devel)",
		},
		Metrics: []metricdata.Metrics{
			{
				Name:        "rpc.server.panics",
				Description: "Measures the number of panics per RPC.",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(
								semconv.RPCService("acme.foo.v1.FooService"),
								semconv.RPCMethod("Bar"),
							),
							Value: 1,
						},
					},
				},
			},
		},
	}

	got, err := collect(t.Context())
	require.NoError(t, err)

	if assert.Len(t, got.ScopeMetrics, 1) {
		metricdatatest.AssertEqual(
			t,
			want,
			got.ScopeMetrics[0],
			metricdatatest.IgnoreTimestamp(),
		)
	}
}
