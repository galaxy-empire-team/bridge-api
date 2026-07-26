package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) CreateMissionEvent(ctx context.Context, missionEvent models.MissionEvent) (uint64, error) {
	const createEventQuery = `
		INSERT INTO session_beta.event_missions (
			mission_id,
			user_id,
			planet_from_x,
			planet_from_y,
			planet_from_z,
			planet_to_x, 
			planet_to_y, 
			planet_to_z, 
			fleet,
			cargo,
			is_returning,
			started_at,
			finished_at
		) VALUES (
			$1,    -- mission_id
			$2,    -- user_id
			$3,    -- planet_from_x
			$4,    -- planet_from_y
			$5,    -- planet_from_z
			$6,    -- planet_to_x
			$7,    -- planet_to_y
			$8,    -- planet_to_z
			$9,    -- fleet
			$10,    -- cargo
			$11,   -- is_returning
			$12,   -- started_at
			$13	   -- finished_at
		)  
		RETURNING id
	`

	fleetJson, err := json.Marshal(toFleetUnits(missionEvent.Fleet))
	if err != nil {
		return 0, fmt.Errorf("json.Marshal(): %w", err)
	}

	cargoJson, err := json.Marshal(toResources(missionEvent.Cargo))
	if err != nil {
		return 0, fmt.Errorf("json.Marshal(): %w", err)
	}

	var eventID uint64
	err = s.DB.QueryRow(ctx, createEventQuery,
		missionEvent.MissionID,
		missionEvent.UserID,
		missionEvent.PlanetFrom.X,
		missionEvent.PlanetFrom.Y,
		missionEvent.PlanetFrom.Z,
		missionEvent.PlanetTo.X,
		missionEvent.PlanetTo.Y,
		missionEvent.PlanetTo.Z,
		fleetJson,
		cargoJson,
		missionEvent.IsReturning,
		missionEvent.StartedAt,
		missionEvent.FinishedAt,
	).Scan(&eventID)
	if err != nil {
		return 0, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return eventID, nil
}
