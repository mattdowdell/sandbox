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
	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect/mockauthnv1connect"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/mockconnect"
)

const (
	testToken      = "an.example.jwt" //#nosec:G101 // not a hardcoded credential
	testAuthHeader = "Bearer " + testToken
)

var mockCtx = mock.AnythingOfType("*context.valueCtx")

// Request implements connect.AnyRequest, but delegates some methods to a generated mock.
//
// The generated mock cannot be used directly because connect.AnyRequest has private methods.
// Additionally, connect.Request cannot be used directly because fields needed by the test are
// private.
type Request struct {
	*connect.Request[healthv1.HealthCheckRequest]

	mock *mockconnect.MockAnyRequest
}

func (r *Request) Spec() connect.Spec {
	return r.mock.Spec()
}

func Test_Interceptor_WrapUnary_Client(t *testing.T) {
	// arrange
	client := mockauthnv1connect.NewMockAuthnServiceClient(t)

	interceptor := authn.New(client)

	inner := mockconnect.NewMockAnyRequest(t)
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
	tests := map[string]struct {
		authn    string
		isClient bool
		options  []authn.Option
	}{
		"uppercase bearer scheme": {
			authn: "BEARER " + testToken,
		},
		"lowercase bearer scheme": {
			authn: "bearer " + testToken,
		},
		"capitalised bearer scheme": {
			authn: testAuthHeader,
		},
		"no authorization with ignore": {
			options: []authn.Option{authn.WithIgnoreService("grpc.health.v1.Health")},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			client := mockauthnv1connect.NewMockAuthnServiceClient(t)
			client.
				EXPECT().
				Authenticate(
					t.Context(),
					connect.NewRequest(&authnv1.AuthenticateRequest{
						Token: testToken,
					}),
				).
				Return(
					connect.NewResponse(&authnv1.AuthenticateResponse{
						Subject: "example",
					}),
					nil,
				).
				Maybe()

			interceptor := authn.New(client, tt.options...)

			inner := mockconnect.NewMockAnyRequest(t)
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
			req.Header().Set("Authorization", tt.authn)

			expected := connect.NewResponse(&healthv1.HealthCheckResponse{})

			next := NewUnaryFunc(t)
			next.
				EXPECT().
				Execute(mockCtx, req).
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
	tests := map[string]struct {
		authn  string
		client func(*testing.T) authnv1connect.AuthnServiceClient
		want   string
	}{
		"missing authorization": {
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()
				return mockauthnv1connect.NewMockAuthnServiceClient(t)
			},
			want: "unauthenticated: invalid or missing authorization",
		},
		"incorrect authorization scheme": {
			authn: "Basic " + testToken,
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()
				return mockauthnv1connect.NewMockAuthnServiceClient(t)
			},
			want: "unauthenticated: invalid or missing authorization",
		},
		"unauthenticated": {
			authn: testAuthHeader,
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewMockAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: testToken,
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
		"unavailable": {
			authn: testAuthHeader,
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewMockAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: testToken,
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
		"internal": {
			authn: testAuthHeader,
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewMockAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: testToken,
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			client := tt.client(t)
			interceptor := authn.New(client)

			req := connect.NewRequest(&healthv1.HealthCheckRequest{})
			req.Header().Set("Authorization", tt.authn)

			next := NewUnaryFunc(t)
			fn := interceptor.WrapUnary(next.Execute)

			// act
			resp, err := fn(t.Context(), req)

			// assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_Interceptor_WrapStreamingHandler_Success(t *testing.T) {
	// arrange
	client := mockauthnv1connect.NewMockAuthnServiceClient(t)
	client.
		EXPECT().
		Authenticate(
			t.Context(),
			connect.NewRequest(&authnv1.AuthenticateRequest{
				Token: testToken,
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
	headers.Set("Authorization", testAuthHeader)

	conn := mockconnect.NewMockStreamingHandlerConn(t)
	conn.EXPECT().Spec().Return(connect.Spec{}).Once()
	conn.EXPECT().RequestHeader().Return(headers).Once()

	next := NewStreamingHandlerFunc(t)
	next.EXPECT().Execute(mockCtx, conn).Return(nil).Once()

	fn := interceptor.WrapStreamingHandler(next.Execute)

	// act
	err := fn(t.Context(), conn)

	// assert
	assert.NoError(t, err)
}

func Test_Interceptor_WrapStreamingHandler_Error(t *testing.T) {
	tests := map[string]struct {
		authn  string
		client func(*testing.T) authnv1connect.AuthnServiceClient
		want   string
	}{
		"missing auth": {
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()
				return mockauthnv1connect.NewMockAuthnServiceClient(t)
			},
			want: "unauthenticated: invalid or missing authorization",
		},
		"auth failed": {
			authn: testAuthHeader,
			client: func(t *testing.T) authnv1connect.AuthnServiceClient {
				t.Helper()

				c := mockauthnv1connect.NewMockAuthnServiceClient(t)
				c.
					EXPECT().
					Authenticate(
						t.Context(),
						connect.NewRequest(&authnv1.AuthenticateRequest{
							Token: testToken,
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			client := tt.client(t)
			interceptor := authn.New(client)

			headers := http.Header{}
			headers.Set("Authorization", tt.authn)

			conn := mockconnect.NewMockStreamingHandlerConn(t)
			conn.EXPECT().Spec().Return(connect.Spec{}).Once()
			conn.EXPECT().RequestHeader().Return(headers).Once()

			next := NewStreamingHandlerFunc(t)
			fn := interceptor.WrapStreamingHandler(next.Execute)

			// act
			err := fn(t.Context(), conn)

			// assert
			assert.EqualError(t, err, tt.want)
		})
	}
}
