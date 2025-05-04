package authn

import (
	"maps"
)

// ...
type Option interface {
	apply(*Interceptor)
}

type ignoreServiceOpt struct {
	services map[string]struct{}
}

// ...
func WithIgnoreService(services ...string) Option {
	return &ignoreServiceOpt{
		services: sliceToMap(services),
	}
}

func (o *ignoreServiceOpt) apply(i *Interceptor) {
	maps.Copy(i.ignores, o.services)
}
