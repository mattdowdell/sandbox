package examplerpc

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid/v5"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// ...
type ResourceFacade interface {
	Create(context.Context, *slog.Logger, *entities.Resource) (*entities.Resource, error)
	Get(context.Context, *slog.Logger, uuid.UUID) (*entities.Resource, error)
	List(context.Context, *slog.Logger, repositories.Pager) (*repositories.Paged[*entities.Resource], error)
	Update(context.Context, *slog.Logger, *entities.Resource) (*entities.Resource, error)
	Delete(context.Context, *slog.Logger, uuid.UUID) error
}

// ...
type AuditEventFacade interface {
	List(context.Context, *slog.Logger) ([]*entities.AuditEvent, error)
	Watch(context.Context, *slog.Logger) <-chan *entities.AuditEvent
}
