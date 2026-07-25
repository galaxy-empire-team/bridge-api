package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PlanetStorage) GetMistLaunches(ctx context.Context, userID uuid.UUID) (uint64, error) {
	const getMistLaunchesQuery = `
		SELECT
			count_left
		FROM session_beta.user_mist_launches
		WHERE user_id = $1;
	`

	var count uint64
	err := r.DB.QueryRow(ctx, getMistLaunchesQuery, userID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return count, nil
}
