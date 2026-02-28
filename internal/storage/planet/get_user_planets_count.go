package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *PlanetStorage) GetUserPlanetsCount(ctx context.Context, userID uuid.UUID) (uint8, error) {
	const getAllUserPlanetsQuery = `
		SELECT count(id)
		FROM session_beta.planets p
		WHERE p.user_id = $1;
	`

	var count uint8
	err := r.DB.QueryRow(ctx, getAllUserPlanetsQuery, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("DB.QueryRow().Scan(): %w", err)
	}

	return count, nil
}
