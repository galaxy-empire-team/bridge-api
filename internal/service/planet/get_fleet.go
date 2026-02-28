package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetFleet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) ([]models.PlanetFleetUnitCount, error) {
	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetID)
	if err != nil {
		return nil, fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return nil, models.ErrPlanetDoesNotBelongToUser
	}

	fleet, err := s.planetStorage.GetFleetForUpdate(ctx, planetID)
	if err != nil {
		return nil, fmt.Errorf("planetStorage.GetFleetForUpdate(): %w", err)
	}

	return fleet, nil
}
