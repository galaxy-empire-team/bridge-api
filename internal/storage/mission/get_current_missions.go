package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) GetCurrentUserMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error) {
	const getCurrentMissionsQuery = `
		SELECT
			ev.mission_type,
			p.x,
			p.y,
			p.z,
			ev.planet_to_x,
			ev.planet_to_y,
			ev.planet_to_z,
			ev.is_returning,
			ev.started_at,
			ev.finished_at
		FROM session_beta.mission_events ev
		JOIN session_beta.planets p ON
			ev.planet_from = p.id
		WHERE ev.user_id = $1;
	`

	rows, err := s.DB.Query(ctx, getCurrentMissionsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}

	var missions []models.UserMission
	for rows.Next() {
		var mission models.UserMission
		err = rows.Scan(
			&mission.Type,
			&mission.PlanetFrom.X,
			&mission.PlanetFrom.Y,
			&mission.PlanetFrom.Z,
			&mission.PlanetTo.X,
			&mission.PlanetTo.Y,
			&mission.PlanetTo.Z,
			&mission.IsReturning,
			&mission.StartedAt,
			&mission.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		missions = append(missions, mission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return missions, nil
}
