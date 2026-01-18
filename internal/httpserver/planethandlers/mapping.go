package planethandlers

import "initialservice/internal/models"

func toTransportPlanet(p models.Planet) PlanetResponse {
	return PlanetResponse{
		PlanetID: p.ID,
		X:        p.X,
		Y:        p.Y,
		Z:        p.Z,
		HasMoon:   p.HasMoon,
		Resource: PlanetResources{
			Metal:   p.Resources.Metal,
			Crystal: p.Resources.Crystal,
			Gas:     p.Resources.Gas,
		},
	}
}