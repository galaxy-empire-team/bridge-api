package mission

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) calculateTotalResources(fleet []models.FleetUnitCount, cargo models.Resources) (models.Resources, error) {
	var gasStartCost uint64
	for _, unit := range fleet {
		unitStats, err := s.registry.GetFleetUnitStatsByID(unit.ID)
		if err != nil {
			return models.Resources{}, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
		}

		gasStartCost += unitStats.GasStartCost * unit.Count
	}

	totalResources := models.Resources{
		Metal:   cargo.Metal,
		Crystal: cargo.Crystal,
		Gas:     gasStartCost + cargo.Gas,
	}

	return totalResources, nil
}
