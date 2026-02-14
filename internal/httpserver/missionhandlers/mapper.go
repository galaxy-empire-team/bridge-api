package missionhandlers

import (
	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func toCoordinatesModel(coordinates Coordinates) models.Coordinates {
	return models.Coordinates{
		X: coordinates.X,
		Y: coordinates.Y,
		Z: coordinates.Z,
	}
}

func toFleetUnits(fleet []FleetUnitCount) []models.PlanetFleetUnitCount {
	units := make([]models.PlanetFleetUnitCount, 0, len(fleet))

	for _, f := range fleet {
		units = append(units, models.PlanetFleetUnitCount{
			ID:    consts.FleetUnitID(f.ID),
			Count: f.Count,
		})
	}

	return units
}
