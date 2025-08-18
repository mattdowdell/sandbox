package authn_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_NewExtractor(t *testing.T) {
	// arrange

	// act
	extractor := authn.NewExtractor()

	// assert
	assert.NotNil(t, extractor)
}

func Test_Extractor_Extract(t *testing.T) {
	testCases := []struct {
		name string
		have func(context.Context) context.Context
		want []slog.Attr
	}{
		{
			name: "present",
			have: func(ctx context.Context) context.Context {
				return jwtx.SubjectIntoContext(ctx, "example")
			},
			want: []slog.Attr{
				slogx.Subject("example"),
			},
		},
		{
			name: "not present",
			have: func(ctx context.Context) context.Context {
				return ctx
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			extractor := authn.NewExtractor()
			ctx := tc.have(t.Context())

			// act
			got := extractor.Extract(ctx)

			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}
