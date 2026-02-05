package configv1client

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"

	"github.com/mattdowdell/sandbox/gen/config/v1"
	"github.com/mattdowdell/sandbox/gen/config/v1/configv1connect"
	"github.com/mattdowdell/sandbox/tests/utils/interceptors"
)

// ...
type Client struct {
	inner configv1connect.ConfigServiceClient
}

// ...
func New(baseURL string) (*Client, error) {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, err
	}

	inner := configv1connect.NewConfigServiceClient(
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
func (c *Client) GetConfig(ctx context.Context) (map[string]string, error) {
	resp, err := c.inner.GetConfig(
		ctx,
		connect.NewRequest(&configv1.GetConfigRequest{}),
	)
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetConfig(), nil
}

// ...
func (c *Client) GetConfigValue(ctx context.Context, key string) (string, error) {
	resp, err := c.inner.GetConfigValue(
		ctx,
		connect.NewRequest(&configv1.GetConfigValueRequest{
			Key: key,
		}),
	)
	if err != nil {
		return "", err
	}

	return resp.Msg.GetValue(), nil
}
