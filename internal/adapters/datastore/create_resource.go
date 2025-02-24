package datastore

import (
	"context"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (s *Stub) CreateResource(ctx context.Context, resource *entities.Resource) error {
	m := modelhelpers.ResourceFromDomain(resource)

	stmt := table.Resources.
		INSERT(
			table.Resources.ID,
			table.Resources.Name,
			table.Resources.CreatedAt,
			table.Resources.UpdatedAt,
		).
		MODEL(m)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
