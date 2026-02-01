package systemhandlers

import "github.com/galaxy-empire-team/bridge-api/internal/models"

func fromModelSystemPlanets(sp models.SystemPlanets) SystemPlanetsResponse {
	response := SystemPlanetsResponse{
		X:       sp.System.X,
		Y:       sp.System.Y,
		Planets: make(map[uint64]PlanetInfo),
	}

	for _, planet := range sp.Planets {
		response.Planets[planet.Z] = PlanetInfo{
			ID:        planet.ID,
			Type:      planet.Type,
			UserLogin: planet.UserLogin,
			HasMoon:   planet.HasMoon,
		}
	}

	return response
}
