package examplerpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// ...
//
//nolint:iface // similar to ResourceUpdater
type ResourceCreator interface {
	Execute(context.Context, repositories.Resource, *entities.Resource) (*entities.Resource, error)
}

// ...
type ResourceGetter interface {
	Execute(context.Context, repositories.Resource, uuid.UUID) (*entities.Resource, error)
}

// ...
type ResourceLister interface {
	Execute(context.Context, repositories.Resource) ([]*entities.Resource, error)
}

// ...
//
//nolint:iface // similar to ResourceCreator
type ResourceUpdater interface {
	Execute(context.Context, repositories.Resource, *entities.Resource) (*entities.Resource, error)
}

// ...
type ResourceDeleter interface {
	Execute(context.Context, repositories.Resource, uuid.UUID) error
}

// ...
type AuditEventLister interface {
	Execute(context.Context, repositories.AuditEvent) ([]*entities.AuditEvent, error)
}

type AuditEventWatcher interface {
	Execute(context.Context, repositories.AuditEvent) <-chan *entities.AuditEvent
}
