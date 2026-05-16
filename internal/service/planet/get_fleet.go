package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetFleet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Fleet, error) {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return models.Fleet{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	fleet, err := s.planetStorage.GetFleetForUpdate(ctx, planetID)
	if err != nil {
		return models.Fleet{}, fmt.Errorf("planetStorage.GetFleetForUpdate(): %w", err)
	}

	fleetConstruction, err := s.planetStorage.GetCurrentFleetConstruction(ctx, planetID)
	if err != nil {
		return models.Fleet{}, fmt.Errorf("planetStorage.GetCurrentFleetConstruction(): %w", err)
	}

	return models.Fleet{
		Fleet:             fleet,
		FleetConstruction: fleetConstruction,
	}, nil
}
