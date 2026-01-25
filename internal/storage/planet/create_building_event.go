package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *PlanetStorage) CreateBuildingEvent(ctx context.Context, buildEvent models.BuildEvent) error {
	const setBuildLvlQuery = `
		INSERT INTO session_beta.building_events (
			planet_id,
			build_type, 
			started_at,
			finished_at
		) VALUES (
			$1,    -- planet_id
			$2,    -- build_type
			NOW(), -- started_at
			$3	   -- finished_at
		)  
		`

	_, err := s.DB.Exec(ctx, setBuildLvlQuery,
		buildEvent.PlanetID,
		buildEvent.BuildingType,
		buildEvent.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", models.ErrEventIsAlreadyScheduled)
	}

	return nil
}
