package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetPlanet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Planet, error) {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return models.Planet{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	err := s.repository.RecalcResources(ctx, userID, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("recalcResourcesWithUpdatedTime(): %w", err)
	}

	planet, err := s.planetStorage.GetPlanet(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetStorage.GetPlanet(): %w", err)
	}

	return models.Planet{
		ID:          planetID,
		Coordinates: planet.Coordinates,
		Resources:   planet.Resources,
		IsCapitol:   planet.IsCapitol,
		HasMoon:     planet.HasMoon,
		ColonizedAt: planet.ColonizedAt,
	}, nil
}
