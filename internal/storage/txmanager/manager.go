package txmanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/db"
	planetservice "github.com/galaxy-empire-team/bridge-api/internal/service/planet"
)

type TxManager struct {
	pool *db.ConnPool
}

func New(connPool *db.ConnPool) *TxManager {
	return &TxManager{pool: connPool}
}

// ExecPlanetTx implemets methods required by planet service. I decided to copy func for each service
// insted of making factories or use empty interfaces.
func (m *TxManager) ExecPlanetTx(
	ctx context.Context,
	handler func(ctx context.Context, storages planetservice.TxStorages) error,
) error {
	return m.exec(ctx, func(tx pgx.Tx) error {
		return handler(ctx, newStorageSet(tx))
	})
}

func (m *TxManager) exec(
	ctx context.Context,
	handler func(tx pgx.Tx) error,
) (err error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin(): %w", err)
	}

	defer func() {
		rollbackErr := tx.Rollback(ctx)

		if rollbackErr == nil || errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return
		}

		if err != nil {
			err = fmt.Errorf("%w; tx.Rollback(): %w", err, rollbackErr)

			return
		}

		err = fmt.Errorf("tx.Rollback(): %w", rollbackErr)
	}()

	err = handler(tx)
	if err != nil {
		return fmt.Errorf("handler(): %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}
