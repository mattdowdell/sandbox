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

// UTCNow returns the current time in UTC.
//
// This should be used when storing the time. It must not be used for calculating the duration
// between 2 points in time as it lacks any monotonic representation.
func (c *Clock) UTCNow() time.Time {
	return time.Now().UTC()
}

// LocalNow returns the current time in the local timezone.
//
// This should only be used to calculate the durations between 2 points in time via the Since and
// Until methods.
func (c *Clock) LocalNow() time.Time {
	return time.Now()
}

// Since returns the time elapsed since the given value.
func (c *Clock) Since(t time.Time) time.Duration {
	return c.LocalNow().Sub(t)
}

// Until returns the duration until the given value.
func (c *Clock) Until(t time.Time) time.Duration {
	return t.Sub(c.LocalNow())
}
