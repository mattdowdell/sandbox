package slogx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_FromContext(t *testing.T) {
	toAdd := slog.New(slog.DiscardHandler)

	tests := map[string]struct {
		have func(context.Context) context.Context
		want *slog.Logger
	}{
		"empty": {
			have: func(ctx context.Context) context.Context {
				return ctx
			},
			want: slog.Default(),
		},
		"nil": {
			have: func(ctx context.Context) context.Context {
				return slogx.IntoContext(ctx, nil /*logger*/)
			},
			want: slog.Default(),
		},
		"present": {
			have: func(ctx context.Context) context.Context {
				return slogx.IntoContext(ctx, toAdd)
			},
			want: toAdd,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			ctx := tt.have(t.Context())

			// act
			logger := slogx.FromContext(ctx)

			// assert
			assert.Same(t, tt.want, logger)
		})
	}
}
