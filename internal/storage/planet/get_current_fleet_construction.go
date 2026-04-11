package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetCurrentFleetConstruction(ctx context.Context, planetID uuid.UUID) (models.FleetUnitConstructionInfo, error) {
	const getFleetConstructionQuery = `
		SELECT 
			fleet_id,
			count,
			started_at,
			finished_at
		FROM session_beta.event_fleet_constructions
		WHERE planet_id = $1
		FOR UPDATE;
	`

	rows, err := r.DB.Query(ctx, getFleetConstructionQuery, planetID)
	if err != nil {
		return models.FleetUnitConstructionInfo{}, fmt.Errorf("DB.Query(): %w", err)
	}
	defer rows.Close()

	var fleetConstruction models.FleetUnitConstructionInfo
	var startedAt, finishedAt time.Time
	if rows.Next() {
		err = rows.Scan(
			&fleetConstruction.FleetID,
			&fleetConstruction.Count,
			&startedAt,
			&finishedAt,
		)
		if err != nil {
			return models.FleetUnitConstructionInfo{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		fleetConstruction.StartedAt = startedAt.UTC()
		fleetConstruction.FinishedAt = finishedAt.UTC()
	}

	return fleetConstruction, nil
}
