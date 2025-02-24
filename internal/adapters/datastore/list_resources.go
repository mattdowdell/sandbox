package datastore

import (
	"context"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/model"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (s *Stub) ListResources(ctx context.Context) ([]*entities.Resource, error) {
	stmt := table.Resources.
		SELECT(
			table.Resources.ID,
			table.Resources.Name,
			table.Resources.CreatedAt,
			table.Resources.UpdatedAt,
		).
		ORDER_BY(table.Resources.ID.ASC())

	var resources []model.Resources
	if err := stmt.QueryContext(ctx, s.db, &resources); err != nil {
		return nil, err
	}

	return modelhelpers.ResourcesToDomain(resources), nil
}
