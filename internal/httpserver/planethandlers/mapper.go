package planethandlers

import (
	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func toPlanetResponse(p models.Planet) PlanetResponse {
	return PlanetResponse{
		PlanetID: p.ID,
		X:        p.Coordinates.X.ToUint8(),
		Y:        p.Coordinates.Y.ToUint16(),
		Z:        p.Coordinates.Z.ToUint8(),
		HasMoon:  p.HasMoon,
		Resources: PlanetResources{
			Metal:   p.Resources.Metal,
			Crystal: p.Resources.Crystal,
			Gas:     p.Resources.Gas,
		},
		IsCapitol: p.IsCapitol,
		BuildingIDs: lo.Map(p.Buildings, func(b consts.BuildingID, _ int) uint16 {
			return b.ToUint16()
		}),
		BuildingsInProgress: lo.Map(p.BuildingsInProgress, func(b models.BuildingInProgress, _ int) BuildingInProgress {
			return BuildingInProgress{
				BuildingID: b.BuildingID,
				StartedAt:  b.StartedAt,
				FinishedAt: b.FinishedAt,
			}
		}),
	}
}

func toUserPlanetsResponse(planets []models.Planet) UserPlanetsResponse {
	resp := UserPlanetsResponse{
		Planets: make([]ShortPlanet, 0, len(planets)),
	}

	for _, p := range planets {
		resp.Planets = append(resp.Planets, ShortPlanet{
			PlanetID:  p.ID,
			X:         p.Coordinates.X.ToUint8(),
			Y:         p.Coordinates.Y.ToUint16(),
			Z:         p.Coordinates.Z.ToUint8(),
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

func toFleetResponse(fleet []models.PlanetFleetUnitCount) FleetPlanetsResponse {
	resp := FleetPlanetsResponse{
		Fleet: make([]FleetUnitCount, 0, len(fleet)),
	}

	for _, p := range fleet {
		resp.Fleet = append(resp.Fleet, FleetUnitCount{
			ID:    p.ID.ToUint8(),
			Count: p.Count,
		})
	}

	return resp
}
