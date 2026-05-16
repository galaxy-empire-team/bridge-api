package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetPlanetResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Resources, error) {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return models.Resources{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	err := s.repository.RecalcResources(ctx, userID, planetID)
	if err != nil {
		return models.Resources{}, fmt.Errorf("recalcResources(): %w", err)
	}

	resources, err := s.planetStorage.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.Resources{}, fmt.Errorf("planetStorage.GetResourcesForUpdate(): %w", err)
	}

	return resources, nil
}
