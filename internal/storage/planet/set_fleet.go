package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) SetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.PlanetFleetUnitCount) error {
	const setFleetQuery = `
		UPDATE session_beta.planet_fleet
		SET count = $1, updated_at = now()
		WHERE planet_id = $2 AND fleet_id = $3;
	`

	var batch pgx.Batch
	for _, fleetUnit := range fleet {
		batch.Queue(setFleetQuery, fleetUnit.Count, planetID, fleetUnit.ID)
	}

	br := r.DB.SendBatch(ctx, &batch)
	defer br.Close()

	for i := 0; i < len(fleet); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch.Exec(): %w", err)
		}
	}

	return nil
}
