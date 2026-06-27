package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetUserBoosts(ctx context.Context, userID uuid.UUID) ([]models.UserBoost, error) {
	const getUserBoostsQuery = `
		SELECT 
			boost_id,
			count
		FROM session_beta.user_boosts
		WHERE user_id = $1;
	`

	rows, err := r.DB.Query(ctx, getUserBoostsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	var boosts []models.UserBoost
	for rows.Next() {
		var boost models.UserBoost
		err = rows.Scan(&boost.ID, &boost.Count)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		boosts = append(boosts, boost)
	}

	return boosts, nil
}
