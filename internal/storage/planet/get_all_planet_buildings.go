package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *PlanetStorage) GetAllPlanetBuildings(ctx context.Context, planetID uuid.UUID) ([]consts.BuildingID, error) {
	const getAllPlanetBuildingsQuery = `
		SELECT 
			building_id
		FROM session_beta.planet_buildings p
		WHERE p.planet_id = $1;
	`

	rows, err := r.DB.Query(ctx, getAllPlanetBuildingsQuery, planetID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	var buildingIDs []consts.BuildingID
	for rows.Next() {
		var buildingID consts.BuildingID
		err = rows.Scan(&buildingID)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		buildingIDs = append(buildingIDs, buildingID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return buildingIDs, nil
}
