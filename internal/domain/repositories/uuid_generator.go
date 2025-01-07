package repositories

import (
	"github.com/google/uuid"
)

// ...
type UUIDGenerator interface {
	// ...
	NewV7() (uuid.UUID, error)
}
