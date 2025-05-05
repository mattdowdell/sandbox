package authn

import (
	"maps"
)

// Option allows customisation when creating a new Interceptor.
type Option interface {
	apply(*Interceptor)
}

type ignoreServiceOpt struct {
	services map[string]struct{}
}

// WithIgnoreService sets the services that should not be authenticated.
//
// It is suggested that this be relatively limited, such as to the health and reflection services.
func WithIgnoreService(services ...string) Option {
	return &ignoreServiceOpt{
		services: sliceToSet(services),
	}
}

func (o *ignoreServiceOpt) apply(i *Interceptor) {
	maps.Copy(i.ignores, o.services)
}
