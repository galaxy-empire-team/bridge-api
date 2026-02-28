package planet

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) UpgradeBuilding(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) error {
	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
	}

	currentBuildsCount, err := s.planetStorage.GetBuildsInProgressCount(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetBuildsInProgressCount(): %w", err)
	}

	if currentBuildsCount >= consts.MaxBuildingsInProgress {
		return models.ErrTooManyBuildingsInProgress
	}

	planetBuildingIDs, err := s.planetStorage.GetAllPlanetBuildings(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetAllPlanetBuildings(): %w", err)
	}

	if !slices.Contains(planetBuildingIDs, buildingID) {
		return models.ErrBuildingNotFound
	}

	err = s.recalcResources(ctx, userID, planetID)
	if err != nil {
		return fmt.Errorf("recalcResources(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		buildEvent, err := s.generateEventForExistingBuilding(ctx, planetID, buildingID, planetRepo)
		if err != nil {
			return fmt.Errorf("generateEventForExistingBuilding(): %w", err)
		}

		err = planetRepo.SetFinishedBuildingTime(ctx, planetID, buildEvent.BuildingID, buildEvent.FinishedAt)
		if err != nil {
			return fmt.Errorf("planetStorage.SetFinishedBuildingTime(): %w", err)
		}

		err = planetRepo.CreateBuildingEvent(ctx, buildEvent)
		if err != nil {
			return fmt.Errorf("planetStorage.CreateBuildingEvent(): %w", err)
		}

		return nil
	})
}

func (s *Service) generateEventForExistingBuilding(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID, planetRepo TxStorages) (models.BuildEvent, error) {
	updatedAt := time.Now()

	// Calculate resources
	resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
	}

	nextLvlBuildingID, err := s.registry.GetBuildingNextLvlID(buildingID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("registry.GetBuildingNextLvlID(): %w", err)
	}

	nextLvlStats, err := s.registry.GetBuildingStatsByID(nextLvlBuildingID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("registry.GetBuildingStatsByID(): %w", err)
	}

	if resources.Metal < nextLvlStats.MetalCost ||
		resources.Crystal < nextLvlStats.CrystalCost ||
		resources.Gas < nextLvlStats.GasCost {
		return models.BuildEvent{}, models.ErrNotEnoughResources
	}

	leftResources := models.Resources{
		Metal:     resources.Metal - nextLvlStats.MetalCost,
		Crystal:   resources.Crystal - nextLvlStats.CrystalCost,
		Gas:       resources.Gas - nextLvlStats.GasCost,
		UpdatedAt: updatedAt,
	}

	err = planetRepo.SetResources(ctx, planetID, leftResources)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("planetRepo.SetResources(): %w", err)
	}

	buildEvent := models.BuildEvent{
		PlanetID:   planetID,
		BuildingID: buildingID,
		StartedAt:  updatedAt,
		FinishedAt: updatedAt.Add(time.Duration(nextLvlStats.UpgradeTimeS) * time.Second),
	}

	return buildEvent, nil
}
