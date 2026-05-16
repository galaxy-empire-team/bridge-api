package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]models.PlanetIDCapitol, error) {
	const getPlanetIDsQuery = `
		SELECT 
			id, 
			is_capitol
		FROM session_beta.planets p
		WHERE p.user_id = $1;
	`

	rows, err := r.DB.Query(ctx, getPlanetIDsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	var planetIDs []models.PlanetIDCapitol
	for rows.Next() {
		var planetID models.PlanetIDCapitol
		err = rows.Scan(
			&planetID.PlanetID,
			&planetID.IsCapitol,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		planetIDs = append(planetIDs, planetID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return planetIDs, nil
}
