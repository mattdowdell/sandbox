package usecases

import (
	"context"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ListResources provides the business logic for listing multiple resources.
type ListResources struct{}

// NewListResources creates a new ListResources.
func NewListResources() *ListResources {
	return &ListResources{}
}

// Execute lists multiple resources.
//
// Any failure will cause ErrInternal to be returned.
func (u *ListResources) Execute(
	ctx context.Context,
	logger *slog.Logger,
	store repositories.Resource,
) ([]*entities.Resource, error) {
	resources, err := store.ListResources(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to list resources", slogx.Err(err))
		return nil, domain.ErrInternal
	}

	return resources, nil
}
