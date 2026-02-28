package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// GetBuildingStats retrieves mine infromation from building stats table: (level, it's production and costs)
func (s *PlanetStorage) GetBuildingStats(ctx context.Context, BuildingType consts.BuildingType, level uint8) (models.BuildingStats, error) {
	const getMineStatQuery = `
		SELECT 
			building_type, 
			level,
			metal_cost,
			crystal_cost,
			gas_cost,
			production_s,
			bonuses,
			upgrade_time_s
		FROM session_beta.s_buildings
		WHERE building_type = $1 AND level = $2;
	`

	var mineInfo models.BuildingStats
	err := s.DB.QueryRow(ctx, getMineStatQuery, BuildingType, level).Scan(
		&mineInfo.Type,
		&mineInfo.Level,
		&mineInfo.MetalCost,
		&mineInfo.CrystalCost,
		&mineInfo.GasCost,
		&mineInfo.ProductionS,
		&mineInfo.Bonuses,
		&mineInfo.UpgradeTimeS,
	)
	if err != nil {
		return models.BuildingStats{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return mineInfo, nil
}
