package examplerpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
//
//nolint:iface // similar to ResourceUpdater
type ResourceCreator interface {
	Execute(context.Context, *entities.Resource) (*entities.Resource, error)
}

// ...
type ResourceGetter interface {
	Execute(context.Context, uuid.UUID) (*entities.Resource, error)
}

// ...
type ResourceLister interface {
	Execute(context.Context) ([]*entities.Resource, error)
}

// ...
//
//nolint:iface // similar to ResourceCreator
type ResourceUpdater interface {
	Execute(context.Context, *entities.Resource) (*entities.Resource, error)
}

// ...
type ResourceDeleter interface {
	Execute(context.Context, uuid.UUID) error
}

// ...
type AuditEventLister interface {
	Execute(context.Context) ([]*entities.AuditEvent, error)
}

type AuditEventWatcher interface {
	Execute(context.Context) <-chan *entities.AuditEvent
}
