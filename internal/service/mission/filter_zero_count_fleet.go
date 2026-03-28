package mission

import "github.com/galaxy-empire-team/bridge-api/internal/models"

func filterZeroCountFleet(fleet []models.FleetUnitCount) []models.FleetUnitCount {
	filtered := make([]models.FleetUnitCount, 0, len(fleet))

	for _, unit := range fleet {
		if unit.Count > 0 {
			filtered = append(filtered, unit)
		}
	}

	return filtered
}
