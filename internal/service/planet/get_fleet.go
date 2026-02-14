package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetFleet(ctx context.Context, planetID uuid.UUID) ([]models.PlanetFleetUnitCount, error) {
	fleet, err := s.planetStorage.GetFleetForUpdate(ctx, planetID)
	if err != nil {
		return nil, fmt.Errorf("getPlanetFleet(): %w", err)
	}

	return fleet, nil
}
