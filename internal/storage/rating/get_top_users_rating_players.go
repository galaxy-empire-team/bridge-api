package rating

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *RatingStorage) GetTopUsersRatingPlayers(ctx context.Context, version uint32, limit uint32) ([]models.RatingPlayer, error) {
	const getTopUsersRatingPlayersQuery = `
		SELECT
			ru.user_id,
			u.login,
			ru.spent_resources,
			ru.rank,
			ru.previous_rank
		FROM session_beta.r_users ru
		JOIN session_beta.users u ON u.id = ru.user_id
		WHERE ru.version = $1
		ORDER BY ru.rank ASC
		LIMIT $2;
	`

	rows, err := r.DB.Query(ctx, getTopUsersRatingPlayersQuery, version, limit)
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
			&player.SpentResources,
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
