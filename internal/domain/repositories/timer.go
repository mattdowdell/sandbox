package repositories

import (
	"time"
)

// Timer implementations provide the current time and operations based on the current time.
type Timer interface {
	// UTCNow returns the current time in UTC.
	//
	// This should be used when storing the time. It must not be used for calculating the duration
	// between 2 points in time as it lacks any monotonic representation.
	UTCNow() time.Time

	// Now returns the current time in the local timezone.
	//
	// This should only be used to calculate the durations between 2 points in time via the Since
	// and Until methods.
	Now() time.Time

	// Since returns the time elapsed since the given value. It is shorthand for Timer.Now().Sub(t).
	Since(t time.Time) time.Duration

	// Until returns the duration until the given value. It is shorthand for t.Sub(Timer.Now()).
	Until(t time.Time) time.Duration
}
