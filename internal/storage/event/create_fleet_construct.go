package event

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *EventStorage) CreateFleetConstructEvent(ctx context.Context, fleetConstructEvent models.FleetConstructEvent) error {
	const createFleetConstructEventQuery = `
		INSERT INTO session_beta.event_fleet_constructions (
			planet_id,
			fleet_id, 
			count,
			metal_cost,
			crystal_cost,
			gas_cost,
			started_at,
			finished_at
		) VALUES (
			$1,    -- planet_id
			$2,    -- fleet_id
			$3,    -- count
			$4,    -- metal_cost
			$5,    -- crystal_cost
			$6,    -- gas_cost
			$7,    -- started_at
			$8	   -- finished_at
		) ON CONFLICT (planet_id) DO NOTHING;
		`

	cmd, err := s.DB.Exec(ctx, createFleetConstructEventQuery,
		fleetConstructEvent.PlanetID,
		fleetConstructEvent.FleetID,
		fleetConstructEvent.Count,
		fleetConstructEvent.ResourcesCost.Metal,
		fleetConstructEvent.ResourcesCost.Crystal,
		fleetConstructEvent.ResourcesCost.Gas,
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
