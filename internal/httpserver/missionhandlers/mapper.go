package missionhandlers

import "github.com/galaxy-empire-team/bridge-api/internal/models"

func toCoordinatesModel(coordinates Coordinates) models.Coordinates {
	return models.Coordinates{
		X: coordinates.X,
		Y: coordinates.Y,
		Z: coordinates.Z,
	}
}
