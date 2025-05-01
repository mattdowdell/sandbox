package examplerpc

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
type ResourceFacade interface {
	Create(context.Context, *slog.Logger, *entities.Resource) (*entities.Resource, error)
	Get(context.Context, *slog.Logger, uuid.UUID) (*entities.Resource, error)
	List(context.Context, *slog.Logger) ([]*entities.Resource, error)
	Update(context.Context, *slog.Logger, *entities.Resource) (*entities.Resource, error)
	Delete(context.Context, *slog.Logger, uuid.UUID) error
}

// ...
type AuditEventFacade interface {
	List(context.Context, *slog.Logger) ([]*entities.AuditEvent, error)
	Watch(context.Context, *slog.Logger) <-chan *entities.AuditEvent
}
