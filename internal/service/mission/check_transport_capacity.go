package mission

import (
	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) checkTransportCapacity(fleet []models.FleetUnitCount, cargo models.Resources) bool {
	var cargoLimit uint64
	for _, fleetUnit := range fleet {
		fStats, err := s.registry.GetFleetUnitStatsByID(fleetUnit.ID)
		if err != nil {
			return false
		}

		cargoLimit += fStats.CargoCapacity * fleetUnit.Count
	}

	return cargo.Metal+cargo.Crystal+cargo.Gas <= cargoLimit
}
