package examplerpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
type ResourceFacade interface {
	Create(context.Context, *entities.Resource) (*entities.Resource, error)
	Get(context.Context, uuid.UUID) (*entities.Resource, error)
	List(context.Context) ([]*entities.Resource, error)
	Update(context.Context, *entities.Resource) (*entities.Resource, error)
	Delete(context.Context, uuid.UUID) error
}

// ...
type AuditEventFacade interface {
	List(context.Context) ([]*entities.AuditEvent, error)
	Watch(context.Context) <-chan *entities.AuditEvent
}
