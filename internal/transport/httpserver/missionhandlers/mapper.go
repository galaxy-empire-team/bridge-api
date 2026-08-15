package missionhandlers

import (
	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func toCoordinatesModel(coordinates Coordinates) models.Coordinates {
	return models.Coordinates{
		X: consts.PlanetPositionX(coordinates.X),
		Y: consts.PlanetPositionY(coordinates.Y),
		Z: consts.PlanetPositionZ(coordinates.Z),
	}
}

func toResources(res Resources) models.Resources {
	return models.Resources{
		Metal:   res.Metal,
		Crystal: res.Crystal,
		Gas:     res.Gas,
	}
}

func fromCoordinatesModel(coordinates models.Coordinates) Coordinates {
	return Coordinates{
		X: coordinates.X.ToUint8(),
		Y: coordinates.Y.ToUint16(),
		Z: coordinates.Z.ToUint8(),
	}
}

func toFleetUnits(fleet []FleetUnitCount) []models.FleetUnitCount {
	units := make([]models.FleetUnitCount, 0, len(fleet))

	for _, f := range fleet {
		units = append(units, models.FleetUnitCount{
			ID:    consts.FleetUnitID(f.ID),
			Count: f.Count,
		})
	}

	return units
}

func fromUserMission(m models.UserMission) Mission {
	return Mission{
		ID:          m.ID,
		PreviousID:  m.PreviousID,
		MissionID:   m.MissionID,
		PlanetFrom:  fromCoordinatesModel(m.PlanetFrom),
		PlanetTo:    fromCoordinatesModel(m.PlanetTo),
		IsReturning: m.IsReturning,
		IsCancelled: m.IsCancelled,
		StartedAt:   m.StartedAt.UTC(),
		FinishedAt:  m.FinishedAt.UTC(),
	}
}

func fromUserMissions(userMissions models.UserMissions) UserMissionsResponse {
	resp := UserMissionsResponse{
		Missions:     make([]Mission, 0, len(userMissions.Missions)),
		MistLaunches: userMissions.MistLaunches,
	}

	for _, m := range userMissions.Missions {
		resp.Missions = append(resp.Missions, fromUserMission(m))
	}

	return resp
}
