package repositories

import (
	"github.com/gofrs/uuid/v5"
)

// UUIDGenerator implementations generate new UUIDs.
type UUIDGenerator interface {
	// NewV7 creates a new UUID v7.
	NewV7() (uuid.UUID, error)
}
