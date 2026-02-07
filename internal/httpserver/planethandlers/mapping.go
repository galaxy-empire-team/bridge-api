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
		X:        p.Coordinates.X,
		Y:        p.Coordinates.Y,
		Z:        p.Coordinates.Z,
		HasMoon:  p.HasMoon,
		Resources: PlanetResources{
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

func toTransportPlanets(planets []models.Planet) UserPlanetsResponse {
	resp := UserPlanetsResponse{
		Planets: make([]GetShortPlanet, 0, len(planets)),
	}

	for _, p := range planets {
		resp.Planets = append(resp.Planets, GetShortPlanet{
			PlanetID:  p.ID,
			X:         p.Coordinates.X,
			Y:         p.Coordinates.Y,
			Z:         p.Coordinates.Z,
			IsCapitol: p.IsCapitol,
			Resources: PlanetResources{
				Metal:   p.Resources.Metal,
				Crystal: p.Resources.Crystal,
				Gas:     p.Resources.Gas,
			},
			ColonizedAt: p.ColonizedAt,
		})
	}

	return resp
}
