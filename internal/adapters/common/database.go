package common

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Datastore is a unified interface for repositories.
type Datastore interface {
	repositories.AuditEvent
	repositories.Resource
}

// CommitFn commits a transaction.
type CommitFn func() error

// RollbackFn rolls back a transaction. No error should be produced for a committed transaction.
type RollbackFn func() error

// Provider provides access to either a transactional or non-transactional database.
type Provider interface {
	Datastore() Datastore
	BeginTx(context.Context) (Datastore, CommitFn, RollbackFn, error)
}

// TxFunc executes the given function within a transaction, automatically committing or rolling back
// the transaction as necessary.
func TxFunc(
	ctx context.Context,
	provider Provider,
	fn func(Datastore) error,
) error {
	_, _, err := TxValues(ctx, provider, func(ds Datastore) (struct{}, struct{}, error) {
		return struct{}{}, struct{}{}, fn(ds)
	})

	return err
}

// TxValue executes the given function within a transaction, automatically committing or rolling
// back the transaction as necessary. It returns the non-error value returned by fn on success.
func TxValue[T any](
	ctx context.Context,
	provider Provider,
	fn func(Datastore) (T, error),
) (T, error) {
	val, _, err := TxValues(ctx, provider, func(ds Datastore) (T, struct{}, error) {
		val, err := fn(ds)
		return val, struct{}{}, err
	})

	return val, err
}

// TxValues executes the given function within a transaction, automatically committing or rolling
// back the transaction as necessary. It returns the non-error values returned by fn on success.
//
//nolint:gocritic // nothing gained by naming generic return values
func TxValues[T1, T2 any](
	ctx context.Context,
	provider Provider,
	fn func(Datastore) (T1, T2, error),
) (T1, T2, error) {
	conn, commit, rollback, err := provider.BeginTx(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin transaction", slogx.Err(err))

		var t1 T1
		var t2 T2

		return t1, t2, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err := rollback(); err != nil {
			slog.ErrorContext(ctx, "failed to rollback transaction", slogx.Err(err))
		}
	}()

	val1, val2, err := fn(conn)
	if err != nil {
		var t1 T1
		var t2 T2

		return t1, t2, err
	}

	if err := commit(); err != nil {
		slog.ErrorContext(ctx, "failed to commit transaction", slogx.Err(err))

		var t1 T1
		var t2 T2

		return t1, t2, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return val1, val2, nil
}
