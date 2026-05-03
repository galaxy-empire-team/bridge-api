package mission

import (
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) checkTransportCapacity(cargo models.Resources, fleet []models.FleetUnitCount, registry registryProvider) bool {
	var cargoLimit uint64
	for _, fleetUnit := range fleet {
		fStats, err := registry.GetFleetUnitStatsByID(fleetUnit.ID)
		if err != nil {
			return false
		}

		cargoLimit += fStats.CargoCapacity * fleetUnit.Count
	}

	return cargo.Metal + cargo.Crystal + cargo.Gas <= cargoLimit
}
