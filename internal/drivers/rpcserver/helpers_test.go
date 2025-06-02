package rpcserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
)

func Test_SplitProcedure(t *testing.T) {
	testCases := []struct {
		name    string
		have    string
		service string
		method  string
	}{
		{
			name:    "service with method",
			have:    "/acme.foo.v1.FooService/Bar",
			service: "acme.foo.v1.FooService",
			method:  "Bar",
		},
		{
			name:   "method only",
			have:   "/Bar",
			method: "Bar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange

			// act
			service, method := rpcserver.SplitProcedure(tc.have)

			// assert
			assert.Equal(t, tc.service, service)
			assert.Equal(t, tc.method, method)
		})
	}
}

func Test_ProtocolToSystem(t *testing.T) {
	testCases := []struct {
		name string
		have string
		want string
	}{
		{
			name: "grpcweb",
			have: "grpcweb",
			want: "grpc_web",
		},
		{
			name: "connectrpc",
			have: "connect",
			want: "connect_rpc",
		},
		{
			name: "grpc",
			have: "grpc",
			want: "grpc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange

			// act
			got := rpcserver.ProtocolToSystem(tc.have)

			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}
