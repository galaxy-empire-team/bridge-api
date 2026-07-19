package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (cp *ConnPool) ExecTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := cp.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin(): %w", err)
	}

	defer func() {
		if err == nil {
			return
		}

		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			err = fmt.Errorf("%w; tx.Rollback(): %w", err, rollbackErr)
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		return fmt.Errorf("fn(): %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}
