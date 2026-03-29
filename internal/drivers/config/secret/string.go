// Package secret provides helpers for secret configuration values, such as passwords or private
// keys.
package secret

import (
	"encoding"
	"fmt"
)

// Non-allocating compile-time checks for interface compliance.
var (
	_ encoding.TextAppender  = (*String)(nil)
	_ encoding.TextMarshaler = (*String)(nil)
	_ fmt.Stringer           = (*String)(nil)
)

// String behaves identically to a string during unmarshaling, but redacts itself during
// unmarshaling to prevent its value being leaked.
type String string

// AppendText implements [encoding.TextAppender].
func (s String) AppendText(b []byte) ([]byte, error) {
	return append(b, "********"...), nil
}

// MarshalText implements [encoding.TextMarshaler].
func (s String) MarshalText() ([]byte, error) {
	return s.AppendText(nil)
}

// String implements [fmt.Stringer].
func (s String) String() string {
	return string(s)
}
