package repositories

import (
	"github.com/google/uuid"
)

// UUIDGenerator implementations generate new UUIDs.
type UUIDGenerator interface {
	// NewV7 creates a new UUID v7.
	NewV7() (uuid.UUID, error)
}
