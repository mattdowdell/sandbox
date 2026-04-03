// Package timex provides [time] functions via a struct.
//
// This is intended to be paired with an interface for the used methods of the struct to be used in
// dependency injection. From there, a mock can be created for the interface, ensuring unit tests
// can use predictable times as needed.
package timex

import (
	"time"
)

// Clock provides a wrapper around functions in [time] to support dependency injection and mocking
// in unit tests.
type Clock struct{}

// NewClock creates a new Clock.
func NewClock() *Clock {
	return &Clock{}
}

// UTCNow returns the current time in UTC.
//
// This should be used when storing the time. It must not be used for calculating the duration
// between 2 points in time as it lacks any monotonic representation.
func (*Clock) UTCNow() time.Time {
	return time.Now().UTC()
}

// Now returns the current time in the local timezone.
//
// This should only be used to calculate the durations between 2 points in time via the Since and
// Until methods.
func (*Clock) Now() time.Time {
	return time.Now()
}

// Since returns the time elapsed since the given value.
func (c *Clock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// Until returns the duration until the given value.
func (c *Clock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}
