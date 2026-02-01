package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error) {
	planetIDs, err := s.planetStorage.GetUserPlanetIDs(ctx, userID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetUserPlanetIDs(): %w", err)
	}

	var capitolID uuid.UUID
	for _, pid := range planetIDs {
		if pid.IsCapitol {
			capitolID = pid.PlanetID
			break
		}
	}

	if capitolID == uuid.Nil {
		return models.Planet{}, models.ErrCapitolNotFound
	}

	updatedAt := time.Now().UTC()
	err = s.recalcResourcesWithUpdatedTime(ctx, capitolID, updatedAt)
	if err != nil {
		return models.Planet{}, fmt.Errorf("recalcResourcesWithUpdatedTime(): %w", err)
	}

	capitolCoordinates, err := s.planetStorage.GetCoordinates(ctx, capitolID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetCapitol(): %w", err)
	}

	resources, err := s.planetStorage.GetResources(ctx, capitolID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetResources(): %w", err)
	}

	buildings, err := s.getBuildingsInfo(ctx, capitolID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("GetBuildingsInfo(): %w", err)
	}

	return models.Planet{
		ID:          capitolID,
		UserID:      userID,
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
