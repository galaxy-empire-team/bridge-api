package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// GetPlanetMinesProduction retrieves mine information from the target planet's buildings
func (s *PlanetStorage) GetPlanetMinesProduction(ctx context.Context, planetID uuid.UUID) (map[consts.BuildingType]uint64, error) {
	const getMineInfoQuery = `
		SELECT 
			b.building_type,
			b.production_s
		FROM session_beta.planet_buildings pb
		JOIN session_beta.s_buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.building_type = ANY($2);
	`

	rows, err := s.DB.Query(ctx, getMineInfoQuery, planetID, []consts.BuildingType{
		consts.BuildingTypeMetalMine,
		consts.BuildingTypeCrystalMine,
		consts.BuildingTypeGasMine,
	})
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	minesProduction := make(map[consts.BuildingType]uint64)
	for rows.Next() {
		var buildingType consts.BuildingType
		var productionS uint64

		err = rows.Scan(&buildingType, &productionS)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(): %w", err)
		}

		minesProduction[buildingType] = productionS
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return minesProduction, nil
}
