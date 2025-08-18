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

// sliceToSet converts a slices values to a hash set. This enables checking an element is present to
// be executed in constant time.
func sliceToSet[T comparable](s []T) map[T]struct{} {
	m := make(map[T]struct{}, len(s))

	for _, k := range s {
		m[k] = struct{}{}
	}

	return m
}

// bearerToken validates the use of the bearer authorization scheme and returns the contained
// credential for further processing.
//
// The bearer scheme is checked for case-insensitively per [RFC 9110, Section 11.1].
//
// [RFC 9110, Section 11.1]: https://www.rfc-editor.org/rfc/rfc9110.html#name-authentication-scheme
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
