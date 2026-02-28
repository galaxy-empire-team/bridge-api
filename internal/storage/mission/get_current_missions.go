package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) GetCurrentUserMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error) {
	const getCurrentMissionsQuery = `
		SELECT
			ev.mission_id,
			p.x,
			p.y,
			p.z,
			ev.planet_to_x,
			ev.planet_to_y,
			ev.planet_to_z,
			ev.is_returning,
			ev.started_at,
			ev.finished_at
		FROM session_beta.event_missions ev
		JOIN session_beta.planets p ON
			ev.planet_from = p.id
		WHERE ev.user_id = $1;
	`

	rows, err := s.DB.Query(ctx, getCurrentMissionsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	var missions []models.UserMission
	var startedAt, finishedAt time.Time
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
			&startedAt,
			&finishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		mission.StartedAt = startedAt.UTC()
		mission.FinishedAt = finishedAt.UTC()

		missions = append(missions, mission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return missions, nil
}
