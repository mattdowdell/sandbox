package authn

import (
	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
)

// Config provides configuration for an AuthnService client.
type Config struct {
	// The address for the AuthnService server.
	Address string `koanf:"address"`
}

// NewClientFromConfig creates a new AuthnService client using the given configuration.
func NewClientFromConfig(
	client connect.HTTPClient,
	conf Config,
	interceptors ...connect.Interceptor,
) authnv1connect.AuthnServiceClient {
	return NewClient(client, conf.Address, interceptors...)
}

// NewClient creates a new AuthnService client.
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
