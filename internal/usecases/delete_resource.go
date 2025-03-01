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
type DeleteResource struct {}

// ...
func NewDeleteResource() *DeleteResource {
	return &DeleteResource{}
}

// ...
func (u *DeleteResource) Execute(
	ctx context.Context,
	store repositories.Resource,
	id uuid.UUID,
) error {
	if err := store.DeleteResource(ctx, id); err != nil {
		slog.ErrorContext(ctx, "failed to delete resource", slogx.Err(err))
		return errors.New("internal error")
	}

	return nil
}
