package datastore

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/public/table"
	"github.com/mattdowdell/sandbox/internal/domain"
)

// ...
func (d *Datastore) DeleteResource(ctx context.Context, id uuid.UUID) error {
	stmt := table.Resources.
		DELETE().
		WHERE(table.Resources.ID.EQ(postgres.UUID(id)))

	result, err := stmt.ExecContext(ctx, d.db)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	switch count {
	case 0:
		return fmt.Errorf("%w: %s", domain.ErrNotFound, id)

	case 1:
		return nil

	default:
		return fmt.Errorf("%w: too many deletes: %d", domain.ErrInternal, count)
	}
}
