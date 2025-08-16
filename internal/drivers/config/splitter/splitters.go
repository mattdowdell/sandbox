// ...
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

// ...
type Comma []string

// ...
func (c *Comma) UnmarshalText(text []byte) error {
	*c = strings.Split(string(text), ",")
	return nil
}

// ...
func (c *Comma) Unwrap() []string {
	return []string(*c)
}

// ...
type Space []string

// ...
func (s *Space) UnmarshalText(text []byte) error {
	*s = strings.Split(string(text), " ")
	return nil
}

// ...
func (s *Space) Unwrap() []string {
	return []string(*s)
}
