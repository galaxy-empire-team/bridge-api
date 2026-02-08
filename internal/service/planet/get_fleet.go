package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetFleet(ctx context.Context, planetID uuid.UUID) ([]models.PlanetFleetUnit, error) {
	fleet, err := s.planetStorage.GetFleetCount(ctx, planetID)
	if err != nil {
		return nil, fmt.Errorf("getPlanetFleet(): %w", err)
	}

	if len(fleet) == 0 {
		return nil, models.ErrFleetNotFound
	}

	result := make([]models.PlanetFleetUnit, 0, len(fleet))
	for _, fleetUnit := range fleet {
		unitStat, err := s.registry.GetFleetUnitStatsByID(fleetUnit.ID)
		if err != nil {
			return nil, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
		}

		result = append(result, models.PlanetFleetUnit{
			ID: fleetUnit.ID,
			Stats: models.PlanetFleetStats{
				Type: unitStat.Type,
			},
			Count: fleetUnit.Count,
		})
	}

	return result, nil
}
