package datastore

import (
	"context"
	"fmt"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/schema/model"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/schema/table"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// ...
func (d *Datastore) ListResources(
	ctx context.Context,
	pager repositories.Pager,
) (*repositories.Paged[*entities.Resource], error) {
	stmt := table.Resources.
		SELECT(
			table.Resources.ID,
			table.Resources.Name,
			table.Resources.CreatedAt,
			table.Resources.UpdatedAt,
		).
		ORDER_BY(table.Resources.ID.ASC()).
		LIMIT(int64(pager.Limit))

	var resources []*model.Resources
	if err := stmt.QueryContext(ctx, d.db, &resources); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return &repositories.Paged[*entities.Resource]{
		Items: modelhelpers.ResourcesToDomain(resources),
	}, nil
}
