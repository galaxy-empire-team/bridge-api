package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.PlanetFleetUnitCount, error) {
	const getFleetQuery = `
		SELECT 
			fleet_id,
			count 
		FROM session_beta.planet_fleet 
		WHERE planet_id = $1
		FOR UPDATE;
	`

	rows, err := r.DB.Query(ctx, getFleetQuery, planetID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}

	var fleet []models.PlanetFleetUnitCount
	for rows.Next() {
		var fleetUnit models.PlanetFleetUnitCount
		err = rows.Scan(
			&fleetUnit.ID,
			&fleetUnit.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		fleet = append(fleet, fleetUnit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return fleet, nil
}
