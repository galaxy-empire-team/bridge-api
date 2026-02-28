package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) getPlanetByID(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Planet, error) {
	if planetID == uuid.Nil {
		return models.Planet{}, models.ErrCapitolNotFound
	}

	updatedAt := time.Now().UTC()
	err := s.recalcResourcesWithUpdatedTime(ctx, userID, planetID, updatedAt)
	if err != nil {
		return models.Planet{}, fmt.Errorf("recalcResourcesWithUpdatedTime(): %w", err)
	}

	capitolCoordinates, err := s.planetStorage.GetCoordinates(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetStorage.GetCoordinates(): %w", err)
	}

	resources, err := s.planetStorage.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetStorage.GetResourcesForUpdate(): %w", err)
	}

	planetBuildingIDs, err := s.planetStorage.GetAllPlanetBuildings(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("GetAllPlanetBuildings(): %w", err)
	}

	buildingsInProgress, err := s.planetStorage.GetCurrentBuilds(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("GetCurrentBuilds(): %w", err)
	}

	return models.Planet{
		ID:                  planetID,
		Coordinates:         capitolCoordinates,
		Resources:           resources,
		Buildings:           planetBuildingIDs,
		BuildingsInProgress: buildingsInProgress,
		IsCapitol:           true,
		HasMoon:             false,
		UpdatedAt:           updatedAt,
	}, nil
}
