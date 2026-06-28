package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *EventStorage) GetFleetConstructionEventForUpdate(ctx context.Context, planetID uuid.UUID) (models.EventFinishTime, error) {
	const getFleetConstructionEventForUpdateQuery = `
			SELECT 
				id,
				started_at,
				finished_at
			FROM session_beta.event_fleet_constructions
			WHERE planet_id = $1
			FOR UPDATE;
		`

	var fleetConstructionEvent models.EventFinishTime
	err := s.DB.QueryRow(ctx, getFleetConstructionEventForUpdateQuery, planetID).Scan(
		&fleetConstructionEvent.EventID,
		&fleetConstructionEvent.StartedAt,
		&fleetConstructionEvent.FinishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.EventFinishTime{}, models.ErrEventIsNotScheduled
		}

		return models.EventFinishTime{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return fleetConstructionEvent, nil
}
