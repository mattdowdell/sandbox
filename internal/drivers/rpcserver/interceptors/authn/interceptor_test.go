package authn_test

import (
	"errors"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	authnv1 "github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
	"github.com/mattdowdell/sandbox/mocks/external/connectrpc.com/mockconnect"
	"github.com/mattdowdell/sandbox/mocks/gen/authn/v1/mockauthnv1connect"
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
	client := mockauthnv1connect.NewAuthnServiceClient(t)

	interceptor := authn.New(client)

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

	next := NewUnaryFunc(t)
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
			client := mockauthnv1connect.NewAuthnServiceClient(t)
			client.
				EXPECT().
				Authenticate(
					t.Context(),
					connect.NewRequest(&authnv1.AuthenticateRequest{
						Token: "an.example.jwt",
					}),
				).
				Return(
					connect.NewResponse(&authnv1.AuthenticateResponse{
						Subject: "example",
					}),
					nil,
				).
				Maybe()

			interceptor := authn.New(client, tc.options...)

			inner := mockconnect.NewAnyRequest(t)
			inner.
				EXPECT().
				Spec().
				Return(connect.Spec{
					Procedure: "/grpc.health.v1.Health/Check",
				}).
				Twice()

			req := &Request{
				Request: connect.NewRequest(&healthv1.HealthCheckRequest{}),
				mock:    inner,
			}
			req.Header().Set("Authorization", tc.authn)

			expected := connect.NewResponse(&healthv1.HealthCheckResponse{})

			next := NewUnaryFunc(t)
			next.
				EXPECT().
				Execute(mock.AnythingOfType("*context.valueCtx"), req).
				Return(expected, nil).
				Once()

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
		name   string
		authn  string
		client func(*testing.T) authnv1connect.AuthnServiceClient
		want   string
	}{
		{
			name: "missing authorization",
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()
				return mockauthnv1connect.NewAuthnServiceClient(t)
			},
			want: "unauthenticated: invalid or missing authorization",
		},
		{
			name:  "incorrect authorization scheme",
			authn: "Basic an.example.jwt",
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()
				return mockauthnv1connect.NewAuthnServiceClient(t)
			},
			want: "unauthenticated: invalid or missing authorization",
		},
		{
			name:  "unauthenticated",
			authn: "Bearer an.example.jwt",
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: "an.example.jwt",
						}),
					).
					Return(
						nil,
						connect.NewError(connect.CodeUnauthenticated, errors.New("example")),
					).
					Once()

				return c
			},
			want: "unauthenticated: example",
		},
		{
			name:  "unavailable",
			authn: "Bearer an.example.jwt",
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: "an.example.jwt",
						}),
					).
					Return(
						nil,
						connect.NewError(connect.CodeUnavailable, errors.New("example")),
					).
					Once()

				return c
			},
			want: "unavailable: service unavailable",
		},
		{
			name:  "internal",
			authn: "Bearer an.example.jwt",
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: "an.example.jwt",
						}),
					).
					Return(
						nil,
						connect.NewError(connect.CodeInternal, errors.New("example")),
					).
					Once()

				return c
			},
			want: "internal: internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			client := tc.client(t)
			interceptor := authn.New(client)

			req := connect.NewRequest(&healthv1.HealthCheckRequest{})
			req.Header().Set("Authorization", tc.authn)

			next := NewUnaryFunc(t)
			fn := interceptor.WrapUnary(next.Execute)

			// act
			resp, err := fn(t.Context(), req)

			// assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, tc.want)
		})
	}
}

func Test_Interceptor_WrapStreamingHandler_Success(t *testing.T) {
	// arrange
	client := mockauthnv1connect.NewAuthnServiceClient(t)
	client.
		EXPECT().
		Authenticate(
			t.Context(),
			connect.NewRequest(&authnv1.AuthenticateRequest{
				Token: "an.example.jwt",
			}),
		).
		Return(
			connect.NewResponse(&authnv1.AuthenticateResponse{
				Subject: "example",
			}),
			nil,
		).
		Once()

	interceptor := authn.New(client)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer an.example.jwt")

	conn := mockconnect.NewStreamingHandlerConn(t)
	conn.EXPECT().Spec().Return(connect.Spec{}).Once()
	conn.EXPECT().RequestHeader().Return(headers).Once()

	next := NewStreamingHandlerFunc(t)
	next.EXPECT().Execute(mock.AnythingOfType("*context.valueCtx"), conn).Return(nil).Once()

	fn := interceptor.WrapStreamingHandler(next.Execute)

	// act
	err := fn(t.Context(), conn)

	// assert
	assert.NoError(t, err)
}

func Test_Interceptor_WrapStreamingHandler_Error(t *testing.T) {
	testCases := []struct {
		name   string
		authn  string
		client func(*testing.T) authnv1connect.AuthnServiceClient
		want   string
	}{
		{
			name: "missing auth",
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()
				return mockauthnv1connect.NewAuthnServiceClient(t)
			},
			want: "unauthenticated: invalid or missing authorization",
		},
		{
			name:  "auth failed",
			authn: "Bearer an.example.jwt",
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: "an.example.jwt",
						}),
					).
					Return(
						nil,
						connect.NewError(connect.CodeInternal, errors.New("example")),
					).
					Once()

				return c
			},
			want: "internal: internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			client := tc.client(t)
			interceptor := authn.New(client)

			headers := http.Header{}
			headers.Set("Authorization", tc.authn)

			conn := mockconnect.NewStreamingHandlerConn(t)
			conn.EXPECT().Spec().Return(connect.Spec{}).Once()
			conn.EXPECT().RequestHeader().Return(headers).Once()

			next := NewStreamingHandlerFunc(t)
			fn := interceptor.WrapStreamingHandler(next.Execute)

			// act
			err := fn(t.Context(), conn)

			// assert
			assert.EqualError(t, err, tc.want)
		})
	}
}
