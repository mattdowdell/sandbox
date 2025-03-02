package datastore

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/mattdowdell/sandbox/internal/adapters/common"
)

// ...
type Provider struct {
	db *sql.DB
}

// ...
func NewProvider(db *sql.DB) *Provider {
	return &Provider{
		db: db,
	}
}

// ...
func (p *Provider) BeginTx(
	ctx context.Context,
) (common.Datastore, common.CommitFn, common.RollbackFn, error) {
	tx, err := p.db.BeginTx(ctx, nil /*opts*/)
	if err != nil {
		return nil, nil, nil, err
	}

	ds := NewDatastore(tx)

	return ds, tx.Commit, tx.Rollback, nil
}

// ...
func (p *Provider) Datastore() common.Datastore {
	return NewDatastore(p.db)
}

// ...
type Datastore struct {
	db qrm.DB
}

// ...
func NewDatastore(db qrm.DB) *Datastore {
	return &Datastore{
		db: db,
	}
}
