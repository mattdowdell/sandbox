package rpcserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
)

func Test_SplitProcedure(t *testing.T) {
	tests := map[string]struct {
		have    string
		service string
		method  string
	}{
		"service with method": {
			have:    "/acme.foo.v1.FooService/Bar",
			service: "acme.foo.v1.FooService",
			method:  "Bar",
		},
		"method only": {
			have:   "/Bar",
			method: "Bar",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			service, method := rpcserver.SplitProcedure(tt.have)

			// assert
			assert.Equal(t, tt.service, service)
			assert.Equal(t, tt.method, method)
		})
	}
}

func Test_ProtocolToSystem(t *testing.T) {
	tests := map[string]struct {
		have string
		want string
	}{
		"grpcweb": {
			have: "grpcweb",
			want: "grpc_web",
		},
		"connectrpc": {
			have: "connect",
			want: "connect_rpc",
		},
		"grpc": {
			have: "grpc",
			want: "grpc",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			got := rpcserver.ProtocolToSystem(tt.have)

			// assert
			assert.Equal(t, tt.want, got)
		})
	}
}
