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
func ParseOperation(input string) Operation {
	switch input {
	case "created":
		return OperationCreated

	case "updated":
		return OperationModified

	case "deleted":
		return OperationDeleted

	default:
		return 0
	}
}

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
func ParseResourceType(input string) ResourceType {
	switch input {
	case "resource":
		return ResourceTypeResource

	default:
		return 0
	}
}

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
