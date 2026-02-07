package planet

import (
	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func toModelBuildingStats(stats registry.BuildingStats) models.BuildingStats {
	return models.BuildingStats{
		Type:         stats.Type,
		Level:        stats.Level,
		MetalCost:    stats.MetalCost,
		CrystalCost:  stats.CrystalCost,
		GasCost:      stats.GasCost,
		ProductionS:  stats.ProductionS,
		Bonuses:      stats.Bonuses,
		UpgradeTimeS: stats.UpgradeTimeS,
	}
}
