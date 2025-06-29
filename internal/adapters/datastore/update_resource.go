package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgerrcode"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/model"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/table"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
func (d *Datastore) UpdateResource(
	ctx context.Context,
	resource *entities.Resource,
) (*entities.Resource, error) {
	m := modelhelpers.ResourceFromDomain(resource)

	stmt := table.Resources.
		UPDATE(
			table.Resources.Name,
			table.Resources.UpdatedAt,
		).
		WHERE(table.Resources.ID.EQ(postgres.UUID(resource.ID))).
		MODEL(m).
		RETURNING(table.Resources.AllColumns)

	var output model.Resources

	if err := stmt.QueryContext(ctx, d.db, &output); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
		}

		if isPgErr(err, pgerrcode.UniqueViolation) {
			return nil, fmt.Errorf("%w: %w", domain.ErrAlreadyExists, err)
		}

		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return modelhelpers.ResourceToDomain(&output), nil
}
