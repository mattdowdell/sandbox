package datastore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/mattdowdell/sandbox/internal/adapters/txn"
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
) (txn.Datastore, txn.Ender, error) {
	tx, err := p.db.BeginTx(ctx, nil /*opts*/)
	if err != nil {
		return nil, nil, err
	}

	ds := NewDatastore(tx)
	end := newEnder(tx)

	return ds, end, nil
}

// Datastore creates a Datastore without a transaction.
func (p *Provider) Datastore() txn.Datastore {
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

// Ender provides support for ending a transaction.
type Ender struct {
	tx *sql.Tx
}

// newEnder creates a new Ender.
func newEnder(tx *sql.Tx) *Ender {
	return &Ender{
		tx: tx,
	}
}

// Commit is used to commit a transaction.
func (e *Ender) Commit() error {
	return e.tx.Commit()
}

// Rollback is used to rollback a transaction, ignoring any errors caused by rollbacking an already
// committed transaction.
func (e *Ender) Rollback() error {
	if err := e.tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		return err
	}

	return nil
}
