package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (s *Service) GetCapitolPlanet(ctx context.Context, userID uuid.UUID) (models.Planet, error) {
	planetIDs, err := s.planetStorage.GetUserPlanetIDs(ctx, userID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetUserPlanetIDs(): %w", err)
	}

	err = s.recalcResources(ctx, planetIDs[0])
	if err != nil {
		return models.Planet{}, fmt.Errorf("recalcResources(): %w", err)
	}

	capitolPlanet, err := s.planetStorage.GetCapitol(ctx, userID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetCapitol(): %w", err)
	}

	return capitolPlanet, nil
}
