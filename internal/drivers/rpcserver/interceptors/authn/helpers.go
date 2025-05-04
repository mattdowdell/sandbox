package authn

import (
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

var errInvalidAuthorization = errors.New("invalid or missing authorization")

// ...
func sliceToMap[T comparable](s []T) map[T]struct{} {
	m := make(map[T]struct{}, len(s))

	for _, k := range s {
		m[k] = struct{}{}
	}

	return m
}

// ...
//
// Case insensitive prefix match. See RFC 9110 Section 11.1.
func bearerToken(headers http.Header) (string, error) {
	value := headers.Get(authorizationHeader)

	if len(value) < len(bearerPrefix) {
		return "", errInvalidAuthorization
	}

	if !strings.EqualFold(value[:len(bearerPrefix)], bearerPrefix) {
		return "", errInvalidAuthorization
	}

	return value[len(bearerPrefix):], nil
}

// ...
func parseToken(_ string) error {
	return nil
}
