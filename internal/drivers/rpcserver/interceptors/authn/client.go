package authn

import (
	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
)

// ...
type ClientConfig struct {
	// ...
	Address string `koanf:"address"`
}

// ...
func NewClientFromConfig(
	client connect.HTTPClient,
	conf ClientConfig,
	interceptors ...connect.Interceptor,
) authnv1connect.AuthnServiceClient {
	return NewClient(client, conf.Address, interceptors...)
}

// ...
func NewClient(
	client connect.HTTPClient,
	address string,
	interceptors ...connect.Interceptor,
) authnv1connect.AuthnServiceClient {
	return authnv1connect.NewAuthnServiceClient(
		client,
		address,
		connect.WithInterceptors(interceptors...),
	)
}
