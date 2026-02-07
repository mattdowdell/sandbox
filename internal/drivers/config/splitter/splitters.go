// Package splitter provides unmarshalers for Koanf config structs to convert strings to string
// slices using a particular delimiter.
package splitter

import (
	"encoding"
	"fmt"
	"strings"
)

type splitter interface {
	encoding.TextAppender
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	fmt.Stringer

	Unwrap() []string
}

// Non-allocating compile-time checks for interface compliance.
var (
	_ splitter = (*Comma)(nil)
	_ splitter = (*Space)(nil)
)

// Comma unmarshals a string into a slice of strings using a comma as a delimiter.
//
//nolint:recvcheck // need a pointer for UnmarshalText, but nothing else.
type Comma []string

// AppendText implements [encoding.TextAppender].
func (c Comma) AppendText(b []byte) ([]byte, error) {
	return append(b, c.String()...), nil
}

// MarshalText implements [encoding.TextMarshaler].
func (c Comma) MarshalText() ([]byte, error) {
	return c.AppendText(nil)
}

// String implements [fmt.Stringer].
func (c Comma) String() string {
	return strings.Join(c, ",")
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (c *Comma) UnmarshalText(b []byte) error {
	if len(b) > 0 {
		*c = strings.Split(string(b), ",")
	}

	return nil
}

// Unwrap returns the underlying string slice after unmarshaling.
func (c Comma) Unwrap() []string {
	return []string(c)
}

// Comma unmarshals a string into a slice of strings using a space as a delimiter.
//
//nolint:recvcheck // need a pointer for UnmarshalText, but nothing else.
type Space []string

// AppendText implements [encoding.TextAppender].
func (s Space) AppendText(b []byte) ([]byte, error) {
	return append(b, s.String()...), nil
}

// MarshalText implements [encoding.TextMarshaler].
func (s Space) MarshalText() ([]byte, error) {
	return s.AppendText(nil)
}

// String implements [fmt.Stringer].
func (s Space) String() string {
	return strings.Join(s, " ")
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (s *Space) UnmarshalText(text []byte) error {
	if len(text) > 0 {
		*s = strings.Split(string(text), " ")
	}

	return nil
}

// Unwrap returns the underlying string slice after unmarshaling.
func (s Space) Unwrap() []string {
	return []string(s)
}
