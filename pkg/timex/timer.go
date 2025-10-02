package timex

import (
	"time"
)

// Timer provides a wrapper around functions in [time] to support dependency injection and mocking
// in unit tests.
type Timer struct{}

// NewTimer creates a new Timer.
func NewTimer() *Timer {
	return &Timer{}
}

// UTCNow returns the current time in UTC.
//
// This should be used when storing the time. It must not be used for calculating the duration
// between 2 points in time as it lacks any monotonic representation.
func (*Timer) UTCNow() time.Time {
	return time.Now().UTC()
}

// Now returns the current time in the local timezone.
//
// This should only be used to calculate the durations between 2 points in time via the Since and
// Until methods.
func (*Timer) Now() time.Time {
	return time.Now()
}

// Since returns the time elapsed since the given value.
func (t *Timer) Since(value time.Time) time.Duration {
	return t.Now().Sub(value)
}

// Until returns the duration until the given value.
func (t *Timer) Until(value time.Time) time.Duration {
	return value.Sub(t.Now())
}
