package slogt_test

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/pkg/slogt"
)

func Test_New(t *testing.T) {
	// arrange
	var writer strings.Builder

	tb := slogt.NewMockTB(t)
	tb.EXPECT().Helper().Times(3)
	tb.EXPECT().Output().Return(&writer).Once()

	// act
	logger := slogt.New(tb)

	// assert
	assert.NotNil(t, logger)
}

func Test_Text(t *testing.T) {
	// arrange
	var writer strings.Builder

	tb := slogt.NewMockTB(t)
	tb.EXPECT().Helper().Twice()
	tb.EXPECT().Output().Return(&writer).Once()

	// act
	logger := slogt.Text(tb)

	// assert
	assert.NotNil(t, logger)
}

func Test_TextWithOptions_Output(t *testing.T) {
	// arrange
	var writer strings.Builder

	tb := slogt.NewMockTB(t)
	tb.EXPECT().Helper().Once()
	tb.EXPECT().Output().Return(&writer).Once()

	logger := slogt.TextWithOptions(tb, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) > 0 || a.Key != slog.TimeKey {
				return a
			}

			return slog.Attr{}
		},
	})

	// act
	logger.InfoContext(t.Context(), "example")

	// assert
	assert.Equal(t, "level=INFO msg=example\n", writer.String())
}

func Test_JSON(t *testing.T) {
	// arrange
	var writer strings.Builder

	tb := slogt.NewMockTB(t)
	tb.EXPECT().Helper().Twice()
	tb.EXPECT().Output().Return(&writer).Once()

	// act
	logger := slogt.JSON(tb)

	// assert
	assert.NotNil(t, logger)
}

func Test_JSONWithOptions_Output(t *testing.T) {
	// arrange
	var writer strings.Builder

	tb := slogt.NewMockTB(t)
	tb.EXPECT().Helper().Once()
	tb.EXPECT().Output().Return(&writer).Once()

	logger := slogt.JSONWithOptions(tb, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) > 0 || a.Key != slog.TimeKey {
				return a
			}

			return slog.Attr{}
		},
	})

	// act
	logger.InfoContext(t.Context(), "example")

	// assert
	assert.JSONEq(t, `{"level":"INFO","msg":"example"}`, writer.String())
}
