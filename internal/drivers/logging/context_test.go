package logging_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/logging"
)

func Test_FromContext(t *testing.T) {
	toAdd := logging.New(slog.LevelInfo)

	testCases := []struct {
		name string
		have func(context.Context) context.Context
		want *slog.Logger
	}{
		{
			name: "empty",
			have: func(ctx context.Context) context.Context {
				return ctx
			},
			want: slog.Default(),
		},
		{
			name: "nil",
			have: func(ctx context.Context) context.Context {
				return logging.SetContext(ctx, nil /*logger*/)
			},
			want: slog.Default(),
		},
		{
			name: "present",
			have: func(ctx context.Context) context.Context {
				return logging.SetContext(ctx, toAdd)
			},
			want: toAdd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			ctx := tc.have(t.Context())

			// act
			logger := logging.FromContext(ctx)

			// assert
			assert.Same(t, tc.want, logger)
		})
	}
}
