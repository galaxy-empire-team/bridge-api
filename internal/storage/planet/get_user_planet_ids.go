package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *PlanetStorage) GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	const getPlanetIDsQuery = `
		SELECT 
			p.id
		FROM session_beta.planets p
		WHERE p.user_id = $1;
	`

	rows, err := r.DB.Query(ctx, getPlanetIDsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}

	var planetIDs []uuid.UUID
	for rows.Next() {
		var planetID uuid.UUID
		err = rows.Scan(
			&planetID,
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
