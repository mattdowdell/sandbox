package rpcerrors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/rpcerrors"
)

func Test_Errs(t *testing.T) {
	tests := map[string]struct {
		have error
		want string
	}{
		"unimplemented": {
			have: rpcerrors.ErrUnimplemented,
			want: "unimplemented: unimplemented",
		},
		"internal": {
			have: rpcerrors.ErrInternal,
			want: "internal: internal error",
		},
		"unavailable": {
			have: rpcerrors.ErrUnavailable,
			want: "unavailable: service unavailable",
		},
		"unauthenticated": {
			have: rpcerrors.ErrUnauthenticated,
			want: "unauthenticated: missing or invalid authentication",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.EqualError(t, tt.have, tt.want)
		})
	}
}
