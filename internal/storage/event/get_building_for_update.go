package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *EventStorage) GetBuildingEventForUpdate(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID) (models.EventFinishTime, error) {
	const getBuildingEventForUpdateQuery = `
			SELECT 
				id,
				started_at,
				finished_at
			FROM session_beta.event_buildings
			WHERE planet_id = $1 AND building_id = $2
			FOR UPDATE;
		`

	var buildingEvent models.EventFinishTime
	err := s.DB.QueryRow(ctx, getBuildingEventForUpdateQuery, planetID, buildingID).Scan(
		&buildingEvent.EventID,
		&buildingEvent.StartedAt,
		&buildingEvent.FinishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.EventFinishTime{}, models.ErrEventIsNotScheduled
		}

		return models.EventFinishTime{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return buildingEvent, nil
}
