package authenticate

import (
	"context"
	"net/http"

	"connectrpc.com/authn"
)

// ...
type Authenticator struct {
	ignore map[string]struct{}
}

// ...
func NewAuthenticator(ignore []string) *Authenticator {
	m := make(map[string]struct{}, len(ignore))
	for _, key := range ignore {
		m[key] = struct{}{}
	}

	return &Authenticator{
		ignore: ignore,
	}
}

// ...
func (a *Authenticator) Authenticate(ctx context.Context, req *http.Request) (any, error) {
	procedure, _ := authn.InferProcedure(req.URL)

	if _, ok := a.ignore[procedure]; ok {
		span := trace.SpanFromContext(ctx)
		span.AddEvent("authentication skipped")

		logger := slogx.FromContext(ctx)
		logger.DebugContext(ctx, "authentication skipped")

		return nil, nil
	}

	token, ok := authn.BearerToken(req)
	if !ok {
		err := authn.Errorf("invalid authentication")
		err.Meta().Set("WWW-Authenticate", "Bearer") // TODO: set a better value

		return nil, err
	}

	// TODO: parse JWT into claims

	return nil, nil
}
