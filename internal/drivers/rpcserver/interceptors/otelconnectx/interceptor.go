package otelconnectx

import (
	"connectrpc.com/otelconnect"
)

// ...
type Interceptor = otelconnect.Interceptor

// ...
type Config struct {
	// ...
	TrustRemote bool `koanf:"trustremote"`
}

// ...
func (c *Config) toOptions() []otelconnect.Option {
	var opts []otelconnect.Option

	if c.TrustRemote {
		opts = append(opts, otelconnect.WithTrustRemote())
	}

	return opts
}

// ...
func NewFromConfig(conf Config) (*Interceptor, error) {
	return otelconnect.NewInterceptor(conf.toOptions()...)
}
