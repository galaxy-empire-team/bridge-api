package systemhandlers

import (
	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func toSystemPlanetsResponse(sp models.SystemPlanets) SystemPlanetsResponse {
	response := SystemPlanetsResponse{
		X:       sp.System.X,
		Y:       sp.System.Y,
		Planets: make(map[consts.PlanetPositionZ]PlanetInfo),
		NPC: lo.Map(sp.NPC, func(npc models.NPCAttack, _ int) NPC {
			return NPC{
				Z:          npc.Z,
				AttackedAt: npc.AttackedAt.UTC(),
			}
		}),
	}

	for _, planet := range sp.Planets {
		response.Planets[planet.Z] = PlanetInfo{
			ID:        planet.ID,
			Type:      planet.Type,
			UserLogin: planet.UserLogin,
			HasMoon:   planet.HasMoon,
			Debris: Debris{
				Metal:   planet.Debris.Metal,
				Crystal: planet.Debris.Crystal,
			},
			AttackedAt: planet.AttackedAt.UTC(),
		}
	}

	return response
}
