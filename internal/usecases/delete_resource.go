package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type DeleteResource struct {
	store repositories.Resource
}

// ...
func NewDeleteResource(
	store repositories.Resource,
) *DeleteResource {
	return &DeleteResource{
		store: store,
	}
}

// ...
func (u *DeleteResource) Execute(
	ctx context.Context,
	id uuid.UUID,
) error {
	if err := u.store.DeleteResource(ctx, id); err != nil {
		slog.ErrorContext(ctx, "failed to delete resource", slogx.Err(err))
		return errors.New("internal error")
	}

	return nil
}
