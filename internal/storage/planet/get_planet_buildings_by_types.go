package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *PlanetStorage) GetPlanetBuildingsByTypes(ctx context.Context, planetID uuid.UUID, buildingTypes []consts.BuildingType) (map[consts.BuildingType]consts.BuildingID, error) {
	const getPlanetBuildingsByTypesQuery = `
		SELECT 
			sb.id,
			sb.building_type
		FROM session_beta.planet_buildings pb
		JOIN session_beta.s_buildings sb ON sb.id = pb.building_id
		WHERE pb.planet_id = $1 AND sb.building_type = ANY($2);
	`

	rows, err := r.DB.Query(ctx, getPlanetBuildingsByTypesQuery, planetID, buildingTypes)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	buildingsByType := make(map[consts.BuildingType]consts.BuildingID)
	for rows.Next() {
		var buildingID consts.BuildingID
		var buildingType consts.BuildingType
		err = rows.Scan(&buildingID, &buildingType)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		buildingsByType[buildingType] = buildingID
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return buildingsByType, nil
}
