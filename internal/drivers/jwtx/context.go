package jwtx

import "context"

type (
	subCtxKey struct{}
)

// SubjectIntoContext adds a subject claim to the given context.
func SubjectIntoContext(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, subCtxKey{}, subject)
}

// SubjectFromContext extracts a subject claim from the given context, returning false if no value
// was found.
func SubjectFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(subCtxKey{}).(string)
	return val, ok
}
