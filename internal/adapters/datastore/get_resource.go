package datastore

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/model"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (s *Stub) GetResource(ctx context.Context, id uuid.UUID) (*entities.Resource, error) {
	stmt := table.Resources.
		SELECT(
			table.Resources.ID,
			table.Resources.Name,
			table.Resources.CreatedAt,
			table.Resources.UpdatedAt,
		).
		WHERE(table.Resources.ID.EQ(postgres.UUID(id)))

	var resource model.Resources
	if err := stmt.QueryContext(ctx, s.db, &resource); err != nil {
		return nil, err
	}

	return modelhelpers.ResourceToDomain(resource), nil
}
