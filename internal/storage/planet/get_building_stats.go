package planet

import (
	"context"
	"fmt"

	"initialservice/internal/models"
)

// GetBuildingStats retrieves mine infromation from building stats table: (level, it's production and costs)
func (s *PlanetStorage) GetBuildingStats(ctx context.Context, BuildingType models.BuildingType, level uint8) (models.BuildingStats, error) {
	const getMineStatQuery = `
		SELECT 
			type, 
			level,
			metal_cost,
			crystal_cost,
			gas_cost,
			metal_production_s,
			crystal_production_s,
			gas_production_s,
			bonuses,
			upgrade_time_s
		FROM session_beta.buildings
		WHERE type = $1 AND level = $2;
	`

	var mineInfo models.BuildingStats
	err := s.DB.QueryRow(ctx, getMineStatQuery, BuildingType, level).Scan(
		&mineInfo.Type,
		&mineInfo.Level,
		&mineInfo.MetalCost,
		&mineInfo.CrystalCost,
		&mineInfo.GasCost,
		&mineInfo.MetalPerSecond,
		&mineInfo.CrystalPerSecond,
		&mineInfo.GasPerSecond,
		&mineInfo.Bonuses,
		&mineInfo.UpgradeTimeInSeconds,
	)
	if err != nil {
		return models.BuildingStats{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return mineInfo, nil
}
