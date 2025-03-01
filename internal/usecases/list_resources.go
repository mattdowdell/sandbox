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
type ListResources struct {}

// ...
func NewListResources() *ListResources {
	return &ListResources{}
}

// ...
func (u *ListResources) Execute(
	ctx context.Context,
	store repositories.Resource,
) ([]*entities.Resource, error) {
	resources, err := store.ListResources(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list resources", slogx.Err(err))
		return nil, errors.New("internal error")
	}

	return resources, nil
}
