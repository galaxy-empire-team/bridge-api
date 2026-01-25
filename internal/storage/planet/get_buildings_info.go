package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"

	"github.com/google/uuid"
)

// GetBuildsInfo retrieves mine infromation from the target planet's buildings
func (s *PlanetStorage) GetBuildingsInfo(ctx context.Context, planetID uuid.UUID, BuildingTypes []models.BuildingType) (map[models.BuildingType]models.BuildingInfo, error) {
	const getMineInfoQuery = `
		SELECT 
			b.level,
			b.type,
			b.metal_production_s,
			b.crystal_production_s,
			b.gas_production_s,
			b.bonuses,
			pb.updated_at,
			pb.finished_at
		FROM session_beta.planet_buildings pb
		JOIN session_beta.buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.type = ANY($2);
	`

	rows, err := s.DB.Query(ctx, getMineInfoQuery, planetID, BuildingTypes)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}

	buildingsInfo := make(map[models.BuildingType]models.BuildingInfo)
	for rows.Next() {
		var buildingInfo models.BuildingInfo
		err = rows.Scan(
			&buildingInfo.Level,
			&buildingInfo.Type,
			&buildingInfo.MetalPerSecond,
			&buildingInfo.CrystalPerSecond,
			&buildingInfo.GasPerSecond,
			&buildingInfo.Bonuses,
			&buildingInfo.UpdatedAt,
			&buildingInfo.FinishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		buildingsInfo[buildingInfo.Type] = buildingInfo
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return buildingsInfo, nil
}
