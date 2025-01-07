package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type GetResource struct {
	store repositories.Resource
}

// ...
func NewGetResource(
	store repositories.Resource,
) *GetResource {
	return &GetResource{
		store: store,
	}
}

// ...
func (u *GetResource) Execute(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Resource, error) {
	resource, err := u.store.GetResource(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get resource", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	return resource, nil
}
