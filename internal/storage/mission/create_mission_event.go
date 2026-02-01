package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) CreateMissionEvent(ctx context.Context, missionEvent models.MissionEvent) error {
	const createEventQuery = `
		INSERT INTO session_beta.mission_events (
			mission_type,
			user_id,
			planet_from,
			planet_to_x, 
			planet_to_y, 
			planet_to_z, 
			started_at,
			finished_at
		) VALUES (
			$1,    -- mission_type
			$2,    -- user_id
			$3,    -- planet_from
			$4,    -- planet_to_x
			$5,    -- planet_to_y
			$6,    -- planet_to_z
			$7,    -- started_at
			$8	   -- finished_at
		)  
		`

	_, err := s.DB.Exec(ctx, createEventQuery,
		missionEvent.Type,
		missionEvent.UserID,
		missionEvent.PlanetFrom,
		missionEvent.PlanetTo.X,
		missionEvent.PlanetTo.Y,
		missionEvent.PlanetTo.Z,
		missionEvent.StartedAt,
		missionEvent.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
