package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) getPlanetByID(ctx context.Context, planetID uuid.UUID) (models.Planet, error) {
	if planetID == uuid.Nil {
		return models.Planet{}, models.ErrCapitolNotFound
	}

	updatedAt := time.Now().UTC()
	err := s.recalcResourcesWithUpdatedTime(ctx, planetID, updatedAt)
	if err != nil {
		return models.Planet{}, fmt.Errorf("recalcResourcesWithUpdatedTime(): %w", err)
	}

	capitolCoordinates, err := s.planetStorage.GetCoordinates(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetCapitol(): %w", err)
	}

	resources, err := s.planetStorage.GetResources(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetResources(): %w", err)
	}

	buildings, err := s.getBuildingsInfo(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("GetBuildingsInfo(): %w", err)
	}

	return models.Planet{
		ID:          planetID,
		Coordinates: capitolCoordinates,
		Resources:   resources,
		Buildings:   buildings,
		IsCapitol:   true,
		HasMoon:     false,
		UpdatedAt:   updatedAt,
	}, nil
}

func (s *Service) getBuildingsInfo(ctx context.Context, planetID uuid.UUID) (map[consts.BuildingType]models.BuildingInfo, error) {
	buildings, err := s.planetStorage.GetBuildingsInfo(ctx, planetID, consts.GetBuildingTypes())
	if err != nil {
		return nil, fmt.Errorf("planetRepo.GetBuildingsInfo(): %w", err)
	}

	// If mines are not build yet, initialize them with default values
	// TODO: get values from db
	for _, mineType := range consts.GetMineTypes() {
		if _, exists := buildings[mineType]; !exists {
			buildings[mineType] = models.BuildingInfo{
				Type:        mineType,
				Level:       defaultLvl,
				ProductionS: defaultProductionS,
			}
		}
	}

	return buildings, nil
}
