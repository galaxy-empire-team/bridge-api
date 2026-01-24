package planet

import (
	"initialservice/internal/models"
)

func toModelPlanet(p Planet) models.Planet {
	planetLocation := models.Location{
		X: p.X,
		Y: p.Y,
		Z: p.Z,
	}

	return models.Planet{
		ID:          p.ID,
		Location:    planetLocation,
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
