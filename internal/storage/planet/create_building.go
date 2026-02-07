package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// CreateBuilding creates a new building on the target planet.
func (s *PlanetStorage) CreateBuilding(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID) error {
	const createBuildingQuery = `
		INSERT INTO session_beta.planet_buildings (
			planet_id, 
			building_id,
			updated_at,
			finished_at,
			created_at
		) VALUES (
			$1,
			$2,
			Now(),
			Now(),
			Now()
		);
	`

	_, err := s.DB.Exec(
		ctx,
		createBuildingQuery,
		planetID,
		buildingID,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
