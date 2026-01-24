package planet

import (
	"context"
	"fmt"
	"initialservice/internal/models"

	"github.com/google/uuid"
)

func (r *PlanetStorage) GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]models.PlanetIDWithCapitol, error) {
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

	var planetIDs []models.PlanetIDWithCapitol
	for rows.Next() {
		var planetID models.PlanetIDWithCapitol
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
