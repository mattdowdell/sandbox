package datastore

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (d *Datastore) UpdateResource(ctx context.Context, resource *entities.Resource) error {
	m := modelhelpers.ResourceFromDomain(resource)

	stmt := table.Resources.
		UPDATE(
			table.Resources.Name,
			table.Resources.UpdatedAt,
		).
		WHERE(table.Resources.ID.EQ(postgres.UUID(resource.ID))).
		MODEL(m)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return err
	}

	return nil
}
