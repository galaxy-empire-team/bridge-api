package statichandlers

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func toBuildingStatsResponse(buildingStats []registry.BuildingStats) []BuildingStats {
	response := make([]BuildingStats, 0, len(buildingStats))

	for _, stats := range buildingStats {
		response = append(response, BuildingStats{
			ID:          stats.ID,
			Level:       stats.Level,
			Type:        stats.Type,
			MetalCost:   stats.MetalCost,
			CrystalCost: stats.CrystalCost,
			GasCost:     stats.GasCost,
			ProductionS: stats.ProductionS,
			Bonuses: BuildingBonuses{
				FleetBuildSpeed: stats.Bonuses.FleetBuildSpeed,
				ResearchSpeed:   stats.Bonuses.ResearchSpeed,
				BuildSpeed:      stats.Bonuses.BuildSpeed,
			},
			UpgradeTimeS: stats.UpgradeTimeS,
		})
	}

	return response
}
