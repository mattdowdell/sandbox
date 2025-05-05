package authn_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
	"github.com/mattdowdell/sandbox/mocks/external/connectrpc.com/mockconnect"
)

// Request implements connect.AnyRequest, but delegates some methods to a generated mock.
//
// The generated mock cannot be used directly because connect.AnyRequest has private methods.
// Additionally, connect.Request cannot be used directly because fields needed by the test are
// private.
type Request struct {
	*connect.Request[healthv1.HealthCheckRequest]

	mock *mockconnect.AnyRequest
}

func (r *Request) Spec() connect.Spec {
	return r.mock.Spec()
}

func Test_Interceptor_WrapUnary_Client(t *testing.T) {
	// arrange
	interceptor := authn.New()

	inner := mockconnect.NewAnyRequest(t)
	inner.
		EXPECT().
		Spec().
		Return(connect.Spec{
			IsClient: true,
		}).
		Once()

	req := &Request{
		Request: connect.NewRequest(&healthv1.HealthCheckRequest{}),
		mock:    inner,
	}

	expected := connect.NewResponse(&healthv1.HealthCheckResponse{})

	next := mockconnect.NewUnaryFunc(t)
	next.EXPECT().Execute(t.Context(), req).Return(expected, nil).Once()

	fn := interceptor.WrapUnary(next.Execute)

	// act
	resp, err := fn(t.Context(), req)

	// assert
	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Interceptor_WrapUnary_Success(t *testing.T) {
	testCases := []struct {
		name     string
		authn    string
		isClient bool
		options  []authn.Option
	}{
		{
			name:     "is client",
			isClient: true,
		},
		{
			name:  "uppercase bearer scheme",
			authn: "BEARER an.example.jwt",
		},
		{
			name:  "lowercase bearer scheme",
			authn: "bearer an.example.jwt",
		},
		{
			name:  "capitalised bearer scheme",
			authn: "Bearer an.example.jwt",
		},
		{
			name:    "no authorization with ignore",
			options: []authn.Option{authn.WithIgnoreService("grpc.health.v1.Health")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			interceptor := authn.New(tc.options...)

			times := 2
			if tc.isClient {
				times = 1
			}

			inner := mockconnect.NewAnyRequest(t)
			inner.
				EXPECT().
				Spec().
				Return(connect.Spec{
					Procedure: "/grpc.health.v1.Health/Check",
					IsClient:  tc.isClient,
				}).
				Times(times)

			req := &Request{
				Request: connect.NewRequest(&healthv1.HealthCheckRequest{}),
				mock:    inner,
			}
			req.Header().Set("Authorization", tc.authn)

			expected := connect.NewResponse(&healthv1.HealthCheckResponse{})

			next := mockconnect.NewUnaryFunc(t)
			next.EXPECT().Execute(t.Context(), req).Return(expected, nil).Once()

			fn := interceptor.WrapUnary(next.Execute)

			// act
			resp, err := fn(t.Context(), req)

			// assert
			assert.Equal(t, expected, resp)
			assert.NoError(t, err)
		})
	}
}

func Test_Interceptor_WrapUnary_Error(t *testing.T) {
	testCases := []struct {
		name string
		have string
	}{
		{
			name: "missing authorization",
		},
		{
			name: "incorrect authorization scheme",
			have: "Basic an.example.jwt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			interceptor := authn.New()

			req := connect.NewRequest(&healthv1.HealthCheckRequest{})
			req.Header().Set("Authorization", tc.have)

			next := mockconnect.NewUnaryFunc(t)
			fn := interceptor.WrapUnary(next.Execute)

			// act
			resp, err := fn(t.Context(), req)

			// assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, "unauthenticated: invalid or missing authorization")
		})
	}
}

func Test_Interceptor_WrapStreamingHandler_Success(t *testing.T) {
	// arrange
	interceptor := authn.New()

	headers := http.Header{}
	headers.Set("Authorization", "Bearer an.example.jwt")

	conn := mockconnect.NewStreamingHandlerConn(t)
	conn.EXPECT().Spec().Return(connect.Spec{}).Once()
	conn.EXPECT().RequestHeader().Return(headers).Once()

	next := mockconnect.NewStreamingHandlerFunc(t)
	next.EXPECT().Execute(t.Context(), conn).Return(nil).Once()

	fn := interceptor.WrapStreamingHandler(next.Execute)

	// act
	err := fn(t.Context(), conn)

	// assert
	assert.NoError(t, err)
}

func Test_Interceptor_WrapStreamingHandler_Error(t *testing.T) {
	// arrange
	interceptor := authn.New()

	conn := mockconnect.NewStreamingHandlerConn(t)
	conn.EXPECT().Spec().Return(connect.Spec{}).Once()
	conn.EXPECT().RequestHeader().Return(http.Header{}).Once()

	next := mockconnect.NewStreamingHandlerFunc(t)
	fn := interceptor.WrapStreamingHandler(next.Execute)

	// act
	err := fn(t.Context(), conn)

	// assert
	assert.EqualError(t, err, "unauthenticated: invalid or missing authorization")
}
