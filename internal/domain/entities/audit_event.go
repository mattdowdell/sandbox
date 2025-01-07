package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ...
type Operation int

// ...
const (
	OperationCreated Operation = iota + 1
	OperationModified
	OperationDeleted
)

// ...
func (o Operation) String() string {
	switch o {
	case OperationCreated:
		return "created"

	case OperationModified:
		return "updated"

	case OperationDeleted:
		return "deleted"

	default:
		return fmt.Sprintf("(unknown:%d)", o)
	}
}

// ...
type ResourceType int

// ...
const (
	ResourceTypeResource ResourceType = iota + 1
	// other resource types here.
)

// ...
func (r ResourceType) String() string {
	switch r {
	case ResourceTypeResource:
		return "resource"

	default:
		return fmt.Sprintf("(unknown:%d)", r)
	}
}

// ...
type AuditEvent struct {
	ID           uuid.UUID
	Operation    Operation
	CreatedAt    time.Time
	Summary      string
	ResourceID   uuid.UUID
	ResourceType ResourceType
}
