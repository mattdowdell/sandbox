package datastore

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/table"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (d *Datastore) CreateResource(ctx context.Context, resource *entities.Resource) error {
	m := modelhelpers.ResourceFromDomain(resource)

	stmt := table.Resources.
		INSERT(
			table.Resources.ID,
			table.Resources.Name,
			table.Resources.CreatedAt,
			table.Resources.UpdatedAt,
		).
		MODEL(m)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		if isPgErr(err, pgerrcode.UniqueViolation) {
			return fmt.Errorf("%w: %w", domain.ErrAlreadyExists, err)
		}

		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return nil
}
