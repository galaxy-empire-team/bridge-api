package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *MissionStorage) CancelMissionEvent(ctx context.Context, event models.CancelMission) error {
	const deleteMissionEventQuery = `
		UPDATE session_beta.event_missions
		SET 	
			is_returning = $2,
			started_at = $3,
			finished_at = $4
		WHERE id = $1;
	`

	_, err := s.DB.Exec(ctx, deleteMissionEventQuery, event.ID, event.IsReturning, event.StartedAt, event.FinishedAt)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}
	return nil
}
