package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// SetFinishedBuildingTime updates updatedAt and finishedAt times. This denormalization
// is needed to optimize common planet queries.
func (r *PlanetStorage) SetFinishedBuildingTime(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID, finishedAt time.Time) error {
	const setFinishedBuildingQuery = `
		UPDATE session_beta.planet_buildings
		SET 
			updated_at  = Now(),
			finished_at = $3
		WHERE planet_id = $1 
		AND 
			building_id = $2;
	`

	_, err := r.DB.Exec(
		ctx,
		setFinishedBuildingQuery,
		planetID,
		buildingID,
		finishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
