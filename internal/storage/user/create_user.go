package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/galaxy-empire-team/bridge-api/internal/db"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *UserStorage) CreateUser(ctx context.Context, user models.User) error {
	cp, ok := r.DB.(*db.ConnPool)
	if !ok {
		return fmt.Errorf("DB is not a *db.ConnPool")
	}

	err := cp.ExecTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err := r.createUserRow(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("createUserRow(): %w", err)
		}

		err = r.createUserResourcesRow(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("createUserResourcesRow(): %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("ExecTx(): %w", err)
	}

	return nil
}

func (r *UserStorage) createUserRow(ctx context.Context, tx pgx.Tx, user models.User) error {
	const createUserRowQuery = `
		INSERT INTO session_beta.users (id, login, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := tx.Exec(
		ctx,
		createUserRowQuery,
		user.ID,
		user.Login,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case uniqueViolationCode:
				return fmt.Errorf("DB.Pool.QueryRow(): %w", models.ErrUserAlreadyExists)
			}
		}

		return fmt.Errorf("tx.Exec(): %w", err)
	}

	return nil
}

func (r *UserStorage) createUserResourcesRow(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	const createUserResourcesQuery = `
		INSERT INTO session_beta.user_resources (
			user_id,
			matter,
			doreye,
			updated_at
		) VALUES (
			$1,    --- userID
			$2,     --- matter
			$3,     --- doreye
			NOW()  --- updated_at
		);
	`

	_, err := tx.Exec(
		ctx,
		createUserResourcesQuery,
		userID,
		0,
		0,
	)
	if err != nil {
		return fmt.Errorf("tx.Exec(): %w", err)
	}

	return nil
}
