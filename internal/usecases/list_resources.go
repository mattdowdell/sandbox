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
type ListResources struct {
	store repositories.Resource
}

// ...
func NewListResources(
	store repositories.Resource,
) *ListResources {
	return &ListResources{
		store: store,
	}
}

// ...
func (u *ListResources) Execute(
	ctx context.Context,
) ([]*entities.Resource, error) {
	resources, err := u.store.ListResources(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list resources", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	return resources, nil
}
