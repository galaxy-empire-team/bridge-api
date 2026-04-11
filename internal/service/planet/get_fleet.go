package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetFleet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Fleet, error) {
	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetID)
	if err != nil {
		return models.Fleet{}, fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.Fleet{}, models.ErrPlanetDoesNotBelongToUser
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
