package repositories

import (
	"time"
)

// ...
type Clock interface {
	// ...
	Now() time.Time
}
