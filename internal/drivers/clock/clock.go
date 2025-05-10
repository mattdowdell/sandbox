package clock

import (
	"time"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// Non-allocating compile-time check for interface compliance.
var _ repositories.Clock = (*Clock)(nil)

// Clock provides a wrapper around functions in [time] to support dependency injection and mocking
// in unit tests.
type Clock struct{}

// New creates a new Clock.
func New() *Clock {
	return &Clock{}
}

// Now returns the current time in UTC.
func (c *Clock) Now() time.Time {
	return time.Now().UTC()
}

// Since returns the time elapsed since the given value.
func (c *Clock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// Until returns the duration until the given value.
func (c *Clock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}
