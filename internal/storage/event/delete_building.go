package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *EventStorage) DeleteBuildingEvent(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID) error {
	const deleteBuildingEventQuery = `
			DELETE FROM session_beta.event_buildings
			WHERE planet_id = $1 AND building_id = $2;
		`

	cmd, err := s.DB.Exec(ctx, deleteBuildingEventQuery, planetID, buildingID)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrEventIsNotScheduled
	}

	return nil
}
