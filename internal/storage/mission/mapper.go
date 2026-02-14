package mission

import "github.com/galaxy-empire-team/bridge-api/internal/models"

func toFleetUnits(fleet []models.PlanetFleetUnitCount) []FleetUnit {
	units := make([]FleetUnit, 0, len(fleet))

	for _, f := range fleet {
		units = append(units, FleetUnit{
			ID:    f.ID.ToUint8(),
			Count: f.Count,
		})
	}

	return units
}
