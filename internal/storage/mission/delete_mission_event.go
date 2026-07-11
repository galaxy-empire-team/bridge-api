package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) CancelMissionEvent(ctx context.Context, event models.CancelMission) error {
	const updateMissionEventQuery = `
		UPDATE session_beta.event_missions
		SET 	
			planet_to_x = $2,
			planet_to_y = $3,
			planet_to_z = $4,
			is_returning = $5,
			started_at = $6,
			finished_at = $7
		WHERE id = $1;
	`

	_, err := s.DB.Exec(
		ctx,
		updateMissionEventQuery,
		event.ID,
		event.PlanetTo.X,
		event.PlanetTo.Y,
		event.PlanetTo.Z,
		event.IsReturning,
		event.StartedAt,
		event.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}
	return nil
}
