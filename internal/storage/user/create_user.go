package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"initialservice/internal/models"
)

func (r *Repository) CreateUser(ctx context.Context, user User) (User, error) {
	const createUserQuery = `
		INSERT INTO session_beta.users (id, login, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, login;
	`

	var createdUser User

	err := r.DB.Pool.QueryRow(
		ctx,
		createUserQuery,
		user.ID,
		user.Login,
	).Scan(&createdUser.ID, &createdUser.Login)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return User{}, fmt.Errorf("DB.Pool.QueryRow(): %w", models.ErrUserAlreadyExists)
			}
		}

		return User{}, fmt.Errorf("DB.Pool.QueryRow().Scan(): %w", err)
	}

	return createdUser, nil
}
