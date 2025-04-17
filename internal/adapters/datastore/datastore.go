package datastore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/mattdowdell/sandbox/internal/adapters/common"
)

// Provider is used to create a Datastore with an optional transaction.
//
// When using transactions, it is recommended to use TxFunc, TxValue or TxValues in the common
// adapter. This reduces boilerplate by abstracting away the commit and rollback of the transaction.
type Provider struct {
	db *sql.DB
}

// NewProvider creates a new Provider.
func NewProvider(db *sql.DB) *Provider {
	return &Provider{
		db: db,
	}
}

// BeginTx creates a Datastore within a transaction.
func (p *Provider) BeginTx(
	ctx context.Context,
) (common.Datastore, common.CommitFn, common.RollbackFn, error) {
	tx, err := p.db.BeginTx(ctx, nil /*opts*/)
	if err != nil {
		return nil, nil, nil, err
	}

	ds := NewDatastore(tx)

	return ds, tx.Commit, wrapRollback(tx.Rollback), nil
}

// Datastore creates a Datastore without a transaction.
func (p *Provider) Datastore() common.Datastore {
	return NewDatastore(p.db)
}

// Datastore provides various operations for manipulating a database.
type Datastore struct {
	db qrm.DB
}

// NewDatastore creates a new Datastore.
func NewDatastore(db qrm.DB) *Datastore {
	return &Datastore{
		db: db,
	}
}

// wrapRollback suppresses the error from a transaction rollback failure if it is caused by the
// transaction already being committed. Otherwise the error is returned as is.
func wrapRollback(fn common.RollbackFn) common.RollbackFn {
	return func() error {
		if err := fn(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			return err
		}

		return nil
	}
}
