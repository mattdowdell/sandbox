package repositories

import (
	"context"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
type AuditEvent interface {
	// ...
	CreateAuditEvent(context.Context, *entities.AuditEvent) error

	// ...
	ListAuditEvents(context.Context) ([]*entities.AuditEvent, error)

	// ...
	WatchAuditEvents(context.Context, chan<- *entities.AuditEvent) error
}
