package user

import (
	"initialservice/internal/models"
)

func toStoragePlanet(p models.Planet) Planet {
	return Planet{
		ID:          p.ID,
		X:           p.X,
		Y:           p.Y,
		Z:           p.Z,
		HasMoon:     p.HasMoon,
		ColonizedAt: p.ColonizedAt,
		Resources: Resources{
			Metal:     p.Resources.Metal,
			Crystal:   p.Resources.Crystal,
			Gas:       p.Resources.Gas,
			UpdatedAt: p.Resources.UpdatedAt,
		},
	}
}

func toModelPlanet(p Planet) models.Planet {
	return models.Planet{
		ID:          p.ID,
		X:           p.X,
		Y:           p.Y,
		Z:           p.Z,
		HasMoon:     p.HasMoon,
		ColonizedAt: p.ColonizedAt,
		Resources: models.Resources{
			Metal:     p.Resources.Metal,
			Crystal:   p.Resources.Crystal,
			Gas:       p.Resources.Gas,
			UpdatedAt: p.Resources.UpdatedAt,
		},
	}
}

func toStoragePlanetToColonize(p models.Planet) PlanetToColonize {
	return PlanetToColonize{
		ID: p.ID,
		X:  p.X,
		Y:  p.Y,
		Z:  p.Z,
	}
}

func toModelPlanetToColonize(p PlanetToColonize) models.Planet {
	return models.Planet{
		ID: p.ID,
		X:  p.X,
		Y:  p.Y,
		Z:  p.Z,
	}
}