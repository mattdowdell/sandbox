package otelx_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

func Test_Extractor_Extract_Valid(t *testing.T) {
	// arrange
	extractor := otelx.NewExtractor(
		otelx.WithSpanID(true),
		otelx.WithSampled(true),
	)

	ctx := otelx.MustSpan(t.Context(), testTraceID, testSpanID, true /*sampled*/)

	// act
	got := extractor.Extract(ctx)

	// assert
	want := []slog.Attr{
		slog.String("trace_id", testTraceID),
		slog.String("span_id", testSpanID),
		slog.Bool("sampled", true),
	}

	assert.Equal(t, want, got)
}

func Test_Extractor_Extract_Invalid(t *testing.T) {
	// arrange
	extractor := otelx.NewExtractor(
		otelx.WithSpanID(true),
		otelx.WithSampled(true),
	)

	// act
	got := extractor.Extract(t.Context())

	// assert
	assert.Empty(t, got)
}
