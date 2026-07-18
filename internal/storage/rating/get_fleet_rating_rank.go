package rating

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *RatingStorage) GetFleetRatingRank(ctx context.Context, userID uuid.UUID, version uint32) (uint32, error) {
	const getFleetRatingRankQuery = `
		SELECT rank
		FROM session_beta.r_fleet_power
		WHERE user_id = $1 AND version = $2;
	`

	var rank uint32
	err := r.DB.QueryRow(ctx, getFleetRatingRankQuery, userID, version).Scan(&rank)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrUserNotInRating
		}

		return 0, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return rank, nil
}
