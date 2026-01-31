package planethandlers

import "github.com/galaxy-empire-team/bridge-api/internal/models"

func toTransportPlanet(p models.Planet) GetPlanetResponse {
	buildings := make(map[string]BuildingInfo)
	for bType, bInfo := range p.Buildings {
		buildings[string(bType)] = BuildingInfo{
			Level:       bInfo.Level,
			Bonuses:     bInfo.Bonuses,
			ProductionS: bInfo.ProductionS,
			UpdatedAt:   bInfo.UpdatedAt,
			FinishedAt:  bInfo.FinishedAt,
		}
	}

	return GetPlanetResponse{
		PlanetID: p.ID,
		X:        p.Location.X,
		Y:        p.Location.Y,
		Z:        p.Location.Z,
		HasMoon:  p.HasMoon,
		Resource: PlanetResources{
			Metal:   p.Resources.Metal,
			Crystal: p.Resources.Crystal,
			Gas:     p.Resources.Gas,
		},
		IsCapitol: p.IsCapitol,
		Buildings: buildings,
	}
}

func toTransportBuildingStats(bs models.BuildingStats) GetBuildStatsResponse {
	return GetBuildStatsResponse{
		Type:                 string(bs.Type),
		Level:                bs.Level,
		MetalCost:            bs.MetalCost,
		CrystalCost:          bs.CrystalCost,
		GasCost:              bs.GasCost,
		ProductionS:          bs.ProductionS,
		Bonuses:              bs.Bonuses,
		UpgradeTimeInSeconds: bs.UpgradeTimeS,
	}
}
