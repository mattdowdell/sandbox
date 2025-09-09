// Package splitter provides unmarshalers for Koanf config structs to convert strings to string
// slices using a particular delimiter.
package splitter

import (
	"encoding"
	"strings"
)

// Non-allocating compile-time checks for interface compliance.
var (
	_ encoding.TextUnmarshaler = (*Comma)(nil)
	_ encoding.TextUnmarshaler = (*Space)(nil)
)

// Comma unmarshals a string into a slice of strings using a comma as a delimiter.
type Comma []string

// ...
func (c *Comma) UnmarshalText(text []byte) error {
	if len(text) > 0 {
		*c = strings.Split(string(text), ",")
	}

	return nil
}

// Unwrap returns the underlying string slice after unmarshaling.
func (c *Comma) Unwrap() []string {
	return []string(*c)
}

// Comma unmarshals a string into a slice of strings using a space as a delimiter.
type Space []string

// ...
func (s *Space) UnmarshalText(text []byte) error {
	if len(text) > 0 {
		*s = strings.Split(string(text), " ")
	}

	return nil
}

// Unwrap returns the underlying string slice after unmarshaling.
func (s *Space) Unwrap() []string {
	return []string(*s)
}
