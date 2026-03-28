package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *PlanetStorage) CreateFleetConstructEvent(ctx context.Context, fleetConstructEvent models.FleetConstructEvent) error {
	const createFleetConstructEventQuery = `
		INSERT INTO session_beta.event_fleet_constructions (
			planet_id,
			fleet_id, 
			count,
			started_at,
			finished_at
		) VALUES (
			$1,    -- planet_id
			$2,    -- fleet_id
			$3,    -- count
			$4,    -- started_at
			$5	   -- finished_at
		) ON CONFLICT (planet_id) DO NOTHING;
		`

	cmd, err := s.DB.Exec(ctx, createFleetConstructEventQuery,
		fleetConstructEvent.PlanetID,
		fleetConstructEvent.FleetID,
		fleetConstructEvent.Count,
		fleetConstructEvent.StartedAt,
		fleetConstructEvent.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrEventIsAlreadyScheduled
	}

	return nil
}
