package planet

import (
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func fromPlanetModel(p models.Planet) planetToColonize {
	return planetToColonize{
		ID:     p.ID,
		UserID: p.UserID,
		Coordinates: coordinates{
			X: p.Coordinates.X,
			Y: p.Coordinates.Y,
			Z: p.Coordinates.Z,
		},
		HasMoon:   p.HasMoon,
		IsCapitol: p.IsCapitol,
	}
}
