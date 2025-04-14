package logging_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/logging"
)

func Test_Handler_WithAttrs(t *testing.T) {
	// arrange
	handler := logging.Wrap(slog.DiscardHandler, nil /*extractors*/)

	// act
	h := handler.WithAttrs([]slog.Attr{
		slog.Bool("example", true),
	})

	// assert
	assert.IsType(t, (*logging.Handler)(nil), h)
}

func Test_Handler_WithGroup(t *testing.T) {
	// arrange
	handler := logging.Wrap(slog.DiscardHandler, nil /*extractors*/)

	// act
	h := handler.WithGroup("example")

	// assert
	assert.IsType(t, (*logging.Handler)(nil), h)
}
