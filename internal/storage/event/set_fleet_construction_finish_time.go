package event

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *EventStorage) SetFleetConstructionFinishTime(ctx context.Context, fleetConstructionEvent models.EventFinishTime) error {
	const setFleetConstructionFinishTimeQuery = `
		UPDATE session_beta.event_fleet_constructions
		SET finished_at = $1
		WHERE id = $2;
	`

	_, err := s.DB.Exec(ctx, setFleetConstructionFinishTimeQuery, fleetConstructionEvent.FinishedAt, fleetConstructionEvent.EventID)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
