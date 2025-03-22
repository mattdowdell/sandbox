package logging_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

const (
	testTraceID = "0123456789abcdef0123456789abcdef"
	testSpanID  = "0123456789abcdef"
)

func Test_New(t *testing.T) {
	// no arrange necessary

	// act
	logger := logging.New(slog.LevelInfo)

	// assert
	assert.NotNil(t, logger)
}

func Test_NewFromConfig(t *testing.T) {
	// arrange
	conf := logging.Config{
		Level: slog.LevelInfo,
	}

	// act
	logger := logging.NewFromConfig(conf)

	// assert
	assert.NotNil(t, logger)
}

func Test_NewAsDefault(t *testing.T) {
	// no arrange necessary

	// act
	logger := logging.NewAsDefault(slog.LevelInfo, slog.LevelDebug)

	// assert
	assert.NotNil(t, logger)
	assert.Equal(t, slog.LevelInfo, logging.Level.Level())
}

func Test_NewAsDefaultFromConfig(t *testing.T) {
	// arrange
	conf := logging.Config{
		Level:       slog.LevelInfo,
		LegacyLevel: slog.LevelDebug,
	}

	// act
	logger := logging.NewAsDefaultFromConfig(conf)

	// assert
	assert.NotNil(t, logger)
	assert.Equal(t, slog.LevelInfo, logging.Level.Level())
}

func Test_Output(t *testing.T) {
	// arrange
	var b bytes.Buffer

	logger := logging.New(
		slog.LevelInfo,
		logging.WithWriter(&b),
		logging.WithSuppressSource(true),
		logging.WithSuppressTime(true),
	)

	// act
	logger.InfoContext(t.Context(), "example")

	// assert
	assert.JSONEq(t, `{"level":"info","msg":"example"}`, b.String())
}

func Test_Output_WithExtractor(t *testing.T) {
	// arrange
	var b bytes.Buffer
	extractor := otelx.NewExtractor()

	logger := logging.New(
		slog.LevelInfo,
		logging.WithWriter(&b),
		logging.WithSuppressSource(true),
		logging.WithSuppressTime(true),
		logging.WithExtractors(extractor),
	)

	ctx := otelx.MustSpan(t.Context(), testTraceID, testSpanID, true /*sampled*/)

	// act
	logger.InfoContext(ctx, "example")

	// assert
	assert.JSONEq(
		t,
		fmt.Sprintf(`{"level":"info","msg":"example","trace_id":%q}`, testTraceID),
		b.String(),
	)
}
