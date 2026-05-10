package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetBuildings(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Buildings, error) {
	if planetID == uuid.Nil {
		return models.Buildings{}, models.ErrCapitolNotFound
	}

	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetID)
	if err != nil {
		return models.Buildings{}, fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.Buildings{}, models.ErrPlanetDoesNotBelongToUser
	}

	planetBuildingIDs, err := s.planetStorage.GetAllPlanetBuildings(ctx, planetID)
	if err != nil {
		return models.Buildings{}, fmt.Errorf("GetAllPlanetBuildings(): %w", err)
	}

	buildingsInProgress, err := s.planetStorage.GetCurrentBuilds(ctx, planetID)
	if err != nil {
		return models.Buildings{}, fmt.Errorf("GetCurrentBuilds(): %w", err)
	}

	return models.Buildings{
		BuildingIDs:         planetBuildingIDs,
		BuildingsInProgress: buildingsInProgress,
	}, nil
}
