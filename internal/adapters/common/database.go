package common

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

type Datastore interface {
	repositories.AuditEvent
	repositories.Resource
}

type CommitFn func() error

type RollbackFn func() error

// ...
type Provider interface {
	Datastore() Datastore
	BeginTx(context.Context) (Datastore, CommitFn, RollbackFn, error)
}

// ...
func TxFunc(
	ctx context.Context,
	provider Provider,
	fn func(Datastore) error,
) error {
	conn, commit, rollback, err := provider.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := rollback(); err != nil {
			slog.ErrorContext(ctx, "failed to rollback transaction", slogx.Err(err))
		}
	}()

	if err := fn(conn); err != nil {
		return err
	}

	if err := commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ...
func TxValue[T any](
	ctx context.Context,
	provider Provider,
	fn func(Datastore) (T, error),
) (T, error) {
	conn, commit, rollback, err := provider.BeginTx(ctx)
	if err != nil {
		var t T
		return t, err
	}

	defer func() {
		if err := rollback(); err != nil {
			slog.ErrorContext(ctx, "failed to rollback transaction", slogx.Err(err))
		}
	}()

	val, err := fn(conn)
	if err != nil {
		var t T
		return t, err
	}

	if err := commit(); err != nil {
		var t T
		return t, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return val, nil
}

// ...
//
//nolint:gocritic // no value gained from naming generic return values
func TxValues[T1, T2 any](
	ctx context.Context,
	provider Provider,
	fn func(Datastore) (T1, T2, error),
) (T1, T2, error) {
	conn, commit, rollback, err := provider.BeginTx(ctx)
	if err != nil {
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
		var t1 T1
		var t2 T2

		return t1, t2, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return val1, val2, nil
}
