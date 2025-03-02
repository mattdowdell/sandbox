package clock

import (
	"time"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// ...
var _ repositories.Clock = (*Clock)(nil)

// ...
type Clock struct{}

// ...
func New() *Clock {
	return &Clock{}
}

// ...
func (c *Clock) Now() time.Time {
	return time.Now().UTC()
}

// ...
func (c *Clock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// ...
func (c *Clock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}
