package jwtx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
)

func Test_SubjectFromContext(t *testing.T) {
	testCases := []struct {
		name string
		have func(context.Context) context.Context
		want string
	}{
		{
			name: "not found",
			have: func(ctx context.Context) context.Context {
				return ctx
			},
			want: "",
		},
		{
			name: "found",
			have: func(ctx context.Context) context.Context {
				return jwtx.SubjectIntoContext(ctx, "example")
			},
			want: "example",
		},
		{
			name: "wrong type",
			have: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, jwtx.SubCtxKey{}, true)
			},
			want: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			ctx := tc.have(t.Context())

			// act
			got, ok := jwtx.SubjectFromContext(ctx)

			// assert
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.want != "", ok)
		})
	}
}
