package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"initialservice/internal/models"
)

func (r *Repository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	const createUserQuery = `
		INSERT INTO session_beta.users (id, login, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, login;
	`

	dbUser := toStorageUser(user)

	var createdUser User
	err := r.DB.Pool.QueryRow(
		ctx,
		createUserQuery,
		dbUser.ID,
		dbUser.Login,
	).Scan(&createdUser.ID, &createdUser.Login)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return models.User{}, fmt.Errorf("DB.Pool.QueryRow(): %w", models.ErrUserAlreadyExists)
			}
		}

		return models.User{}, fmt.Errorf("DB.Pool.QueryRow().Scan(): %w", err)
	}

	return toModelUser(createdUser), nil
}
