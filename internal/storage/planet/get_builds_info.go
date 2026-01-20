package planet

import (
	"context"
	"fmt"

	"initialservice/internal/models"

	"github.com/google/uuid"
)

// GetBuildsInfo retrieves mine infromation from the target planet's buildings
func (s *PlanetStorage) GetBuildsInfo(ctx context.Context, planetID uuid.UUID, BuildTypes []models.BuildType) (map[models.BuildType]models.PlanetBuildInfo, error) {
	const getMineInfoQuery = `
		SELECT 
			b.level,
			b.type,
			b.metal_production_s,
			b.crystal_production_s,
			b.gas_production_s,
			b.bonuses
		FROM session_beta.planet_buildings pb
		JOIN session_beta.buildings b ON pb.building_id = b.id
		WHERE pb.planet_id = $1 AND b.type = ANY($2);
	`

	rows, err := s.DB.Query(ctx, getMineInfoQuery, planetID, BuildTypes)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}

	minesInfo := make(map[models.BuildType]models.PlanetBuildInfo)
	for rows.Next() {
		var mineInfo models.PlanetBuildInfo
		err = rows.Scan(
			&mineInfo.Level,
			&mineInfo.Type,
			&mineInfo.MetalPerSecond,
			&mineInfo.CrystalPerSecond,
			&mineInfo.GasPerSecond,
			&mineInfo.Bonuses,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		minesInfo[mineInfo.Type] = mineInfo
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return minesInfo, nil
}
