package datastore

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/table"
)

// ...
func (s *Stub) DeleteResource(ctx context.Context, id uuid.UUID) error {
	stmt := table.Resources.
		DELETE().
		WHERE(table.Resources.ID.EQ(postgres.UUID(id)))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
