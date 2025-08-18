package authnv1client

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"

	authnv1 "github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
	"github.com/mattdowdell/sandbox/tests/utils/interceptors"
)

// ...
type Client struct {
	inner authnv1connect.AuthnServiceClient
}

// ...
func New(baseURL string) (*Client, error) {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, err
	}

	inner := authnv1connect.NewAuthnServiceClient(
		http.DefaultClient,
		baseURL,
		connect.WithInterceptors(
			connect.UnaryInterceptorFunc(interceptors.ValidateUnaryInterceptor),
			connect.UnaryInterceptorFunc(interceptors.ScenarioUnaryInterceptor),
			otelInterceptor,
		),
	)

	return &Client{
		inner: inner,
	}, nil
}

// ...
func (c *Client) Login(ctx context.Context, id, secret string) (string, error) {
	resp, err := c.inner.Login(
		ctx,
		connect.NewRequest(&authnv1.LoginRequest{
			Id:     id,
			Secret: secret,
		}),
	)
	if err != nil {
		return "", err
	}

	return resp.Msg.GetAccessToken(), nil
}
