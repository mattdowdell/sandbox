package repositories

import (
	"time"
)

// ...
type Clock interface {
	// Now returns the current time in UTC.
	Now() time.Time

	// Since returns the time elapsed since the given value. It is shorthand for Clock.Now().Sub(t).
	Since(t time.Time) time.Duration

	// Until returns the duration until the given value. It is shorthand for t.Sub(Clock.Now()).
	Until(t time.Time) time.Duration
}
