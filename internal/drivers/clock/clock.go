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

func (c *Clock) Now() time.Time {
	return time.Now().UTC()
}
