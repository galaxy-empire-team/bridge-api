package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) removeFromPlanet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, fleet []models.FleetUnitCount, cargo models.Resources, storage TxStorages) error {
	err := s.repository.RecalcResources(ctx, userID, planetID)
	if err != nil {
		return fmt.Errorf("repository.RecalcResources(): %w", err)
	}

	err = s.removeFleetFromPlanet(ctx, planetID, fleet, storage)
	if err != nil {
		return fmt.Errorf("removeFleetFromPlanet(): %w", err)
	}

	totalResources, err := s.calculateTotalResources(fleet, cargo)
	if err != nil {
		return fmt.Errorf("calculateTotalResources(): %w", err)
	}

	err = s.removeResourcesFromPlanet(ctx, planetID, totalResources, storage)
	if err != nil {
		return fmt.Errorf("removeResourcesFromPlanet(): %w", err)
	}

	return nil
}
