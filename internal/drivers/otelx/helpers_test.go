package otelx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

const (
	testTraceID = "0123456789abcdef0123456789abcdef"
	testSpanID  = "0123456789abcdef"
)

func Test_MustSpan_Success(t *testing.T) {
	assert.NotPanics(t, func() {
		otelx.MustSpan(t.Context(), testTraceID, testSpanID, true /*sampled*/)
	})
}

func Test_MustSpan_Error(t *testing.T) {
	testCases := []struct {
		name    string
		traceID string
		spanID  string
		want    string
	}{
		{
			name:    "invalid trace id",
			traceID: "invalid",
			spanID:  testSpanID,
			want:    "hex encoded trace-id must have length equals to 32",
		},
		{
			name:    "invalid span id",
			traceID: testTraceID,
			spanID:  "invalid",
			want:    "hex encoded span-id must have length equals to 16",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.PanicsWithError(t, tc.want, func() {
				otelx.MustSpan(t.Context(), tc.traceID, tc.spanID, true /*sampled*/)
			})
		})
	}
}
