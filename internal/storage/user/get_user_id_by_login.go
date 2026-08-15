package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *UserStorage) GetUserIDByLogin(ctx context.Context, login string) (uuid.UUID, error) {
	const getUserIDByLoginQuery = `
		SELECT id
		FROM session_beta.users
		WHERE login = $1;
	`

	var userID uuid.UUID
	err := r.DB.QueryRow(ctx, getUserIDByLoginQuery, login).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, models.ErrUserNotFound
		}

		return uuid.Nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return userID, nil
}
