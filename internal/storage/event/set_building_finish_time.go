package event

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *EventStorage) SetBuildingFinishTime(ctx context.Context, buildingEvent models.EventFinishTime) error {
	const setBuildingFinishTimeQuery = `
		UPDATE session_beta.event_buildings
		SET finished_at = $1
		WHERE id = $2;
	`

	_, err := s.DB.Exec(ctx, setBuildingFinishTimeQuery, buildingEvent.FinishedAt, buildingEvent.EventID)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
