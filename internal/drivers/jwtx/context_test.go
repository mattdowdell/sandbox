package jwtx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
)

func Test_SubjectFromContext(t *testing.T) {
	tests := map[string]struct {
		have func(context.Context) context.Context
		want string
	}{
		"not found": {
			have: func(ctx context.Context) context.Context {
				return ctx
			},
			want: "",
		},
		"found": {
			have: func(ctx context.Context) context.Context {
				return jwtx.SubjectIntoContext(ctx, "example")
			},
			want: "example",
		},
		"wrong type": {
			have: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, jwtx.SubCtxKey{}, true)
			},
			want: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			ctx := tt.have(t.Context())

			// act
			got, ok := jwtx.SubjectFromContext(ctx)

			// assert
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want != "", ok)
		})
	}
}
