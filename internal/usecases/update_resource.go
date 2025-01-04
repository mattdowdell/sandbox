package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type UpdateResource struct {
	clock repositories.Clock
	store repositories.Resource
}

// ...
func NewUpdateResource(
	clock repositories.Clock,
	store repositories.Resource,
) *UpdateResource {
	return &UpdateResource{
		clock: clock,
		store: store,
	}
}

// ...
func (u *UpdateResource) Execute(
	ctx context.Context,
	resource *entities.Resource,
) (*entities.Resource, error) {
	resource.Update(u.clock.Now())

	// TODO: handle conflict
	if err := u.store.UpdateResource(ctx, resource); err != nil {
		slog.ErrorContext(ctx, "failed to update resource", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	slog.InfoContext(ctx, "updated resource")

	return resource, nil
}
