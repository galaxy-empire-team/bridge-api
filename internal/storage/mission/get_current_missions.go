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
			id,
			previous_id,
			mission_id,
			planet_from_x,
			planet_from_y,
			planet_from_z,
			planet_to_x,
			planet_to_y,
			planet_to_z,
			is_returning,
			is_cancelled,
			started_at,
			finished_at
		FROM session_beta.event_missions
		WHERE user_id = $1 AND is_finished = false;
	`

	rows, err := s.DB.Query(ctx, getCurrentMissionsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	var (
		missions   []models.UserMission
		previousID *uint64
		startedAt  time.Time
		finishedAt time.Time
	)
	for rows.Next() {
		var mission models.UserMission
		err = rows.Scan(
			&mission.ID,
			&previousID,
			&mission.MissionID,
			&mission.PlanetFrom.X,
			&mission.PlanetFrom.Y,
			&mission.PlanetFrom.Z,
			&mission.PlanetTo.X,
			&mission.PlanetTo.Y,
			&mission.PlanetTo.Z,
			&mission.IsReturning,
			&mission.IsCancelled,
			&startedAt,
			&finishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		mission.StartedAt = startedAt.UTC()
		mission.FinishedAt = finishedAt.UTC()
		if previousID != nil {
			mission.PreviousID = *previousID
		}

		missions = append(missions, mission)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return missions, nil
}
