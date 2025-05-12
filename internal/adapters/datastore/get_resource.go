package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/model"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (d *Datastore) GetResource(ctx context.Context, id uuid.UUID) (*entities.Resource, error) {
	stmt := table.Resources.
		SELECT(
			table.Resources.ID,
			table.Resources.Name,
			table.Resources.CreatedAt,
			table.Resources.UpdatedAt,
		).
		WHERE(table.Resources.ID.EQ(postgres.UUID(id)))

	var resource model.Resources
	if err := stmt.QueryContext(ctx, d.db, &resource); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
		}

		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return modelhelpers.ResourceToDomain(&resource), nil
}
