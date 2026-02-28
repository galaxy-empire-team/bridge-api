package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *PlanetStorage) CreateBuildingEvent(ctx context.Context, buildEvent models.BuildEvent) error {
	const setBuildLvlQuery = `
		INSERT INTO session_beta.event_buildings (
			planet_id,
			building_id, 
			started_at,
			finished_at
		) VALUES (
			$1,    -- planet_id
			$2,    -- building_id
			NOW(), -- started_at
			$3	   -- finished_at
		) ON CONFLICT (planet_id, building_id) DO NOTHING;
		`

	cmd, err := s.DB.Exec(ctx, setBuildLvlQuery,
		buildEvent.PlanetID,
		buildEvent.BuildingID,
		buildEvent.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrEventIsAlreadyScheduled
	}

	return nil
}
