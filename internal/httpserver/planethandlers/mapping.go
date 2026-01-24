package planethandlers

import "initialservice/internal/models"

func toTransportPlanet(p models.Planet) PlanetResponse {
	buildings := make(map[string]BuildingInfo)
	for bType, bInfo := range p.Buildings {
		buildings[string(bType)] = BuildingInfo{
			Level:            bInfo.Level,
			Bonuses:          bInfo.Bonuses,
			MetalPerSecond:   bInfo.MetalPerSecond,
			CrystalPerSecond: bInfo.CrystalPerSecond,
			GasPerSecond:     bInfo.GasPerSecond,
			UpdatedAt:        bInfo.UpdatedAt,
			FinishedAt:       bInfo.FinishedAt,
		}
	}

	return PlanetResponse{
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
