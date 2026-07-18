package rating

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *RatingStorage) GetFleetRatingPlayersByRankOffset(ctx context.Context, version uint32, offset uint32, limit uint32) ([]models.RatingPlayer, error) {
	const getFleetRatingPlayersByRankOffsetQuery = `
		SELECT
			rfp.user_id,
			u.login,
			rfp.fleet_power,
			rfp.rank,
			rfp.previous_rank
		FROM session_beta.r_fleet_power rfp
		JOIN session_beta.users u ON u.id = rfp.user_id
		WHERE rfp.version = $1
		ORDER BY rfp.rank ASC
		OFFSET $2
		LIMIT $3;
	`

	rows, err := r.DB.Query(ctx, getFleetRatingPlayersByRankOffsetQuery, version, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	var players []models.RatingPlayer
	for rows.Next() {
		var player models.RatingPlayer
		err = rows.Scan(
			&player.UserID,
			&player.Login,
			&player.FleetPower,
			&player.Rank,
			&player.PreviousRank,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return players, nil
}
