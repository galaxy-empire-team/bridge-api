package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) AddFleet(ctx context.Context, planetID uuid.UUID, fleetUnit models.FleetUnitCount) error {
	const addFleetQuery = `
		INSERT INTO session_beta.planet_fleet (
			planet_id,
			fleet_id,
			count,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			now()
		)
		ON CONFLICT (planet_id, fleet_id) DO UPDATE SET
			count = session_beta.planet_fleet.count + excluded.count,
			updated_at = now();
	`

	_, err := r.DB.Exec(ctx, addFleetQuery, planetID, fleetUnit.ID, fleetUnit.Count)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
