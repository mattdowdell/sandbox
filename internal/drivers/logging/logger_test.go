package logging_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/logging"
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
