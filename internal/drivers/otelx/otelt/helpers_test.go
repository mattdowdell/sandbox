package otelt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx/otelt"
)

const (
	testTraceID = "0123456789abcdef0123456789abcdef"
	testSpanID  = "0123456789abcdef"
)

func Test_MustSpan_Success(t *testing.T) {
	assert.NotPanics(t, func() {
		otelt.MustSpan(t.Context(), testTraceID, testSpanID, true /*sampled*/)
	})
}

func Test_MustSpan_Error(t *testing.T) {
	tests := map[string]struct {
		traceID string
		spanID  string
		want    string
	}{
		"invalid trace id": {
			traceID: "invalid",
			spanID:  testSpanID,
			want:    "hex encoded trace-id must have length equals to 32",
		},
		"invalid span id": {
			traceID: testTraceID,
			spanID:  "invalid",
			want:    "hex encoded span-id must have length equals to 16",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.PanicsWithError(t, tt.want, func() {
				otelt.MustSpan(t.Context(), tt.traceID, tt.spanID, true /*sampled*/)
			})
		})
	}
}

func Test_MeterProvider(t *testing.T) {
	// arrange

	// act
	provider, collect := otelt.MeterProvider()

	// assert
	assert.NotNil(t, provider)

	if assert.NotNil(t, collect) {
		metrics, err := collect(t.Context())

		assert.Empty(t, metrics.ScopeMetrics)
		assert.NoError(t, err)
	}
}

func Test_TracerProvider(t *testing.T) {
	// arrange

	// act
	provider, collect := otelt.TracerProvider()

	// assert
	assert.NotNil(t, provider)

	if assert.NotNil(t, collect) {
		spans := collect()
		assert.Empty(t, spans)
	}
}
