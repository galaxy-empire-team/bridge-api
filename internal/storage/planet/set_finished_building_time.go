package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (r *PlanetStorage) SetFinishedBuildingTime(ctx context.Context, planetID uuid.UUID, buildingInfo models.BuildingInfo) error {
	const setFinishedBuildingQuery = `
		UPDATE session_beta.planet_buildings
		SET 
			finished_at = $4
		WHERE planet_id = $1 
		AND 
			building_id = (
				SELECT id FROM session_beta.buildings 
				WHERE type = $2 AND level = $3
			);
	`

	finishedBuilding := finishedBuilding{
		Type:       string(buildingInfo.Type),
		Level:      buildingInfo.Level,
		FinishedAt: *buildingInfo.FinishedAt,
	}
	_, err := r.DB.Exec(
		ctx,
		setFinishedBuildingQuery,
		planetID,
		finishedBuilding.Type,
		finishedBuilding.Level,
		finishedBuilding.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
