package planet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) UpgradeBuilding(ctx context.Context, planetID uuid.UUID, BuildingType string) error {
	if !consts.IsValidBuildingType(consts.BuildingType(BuildingType)) {
		return models.ErrBuildTypeInvalid
	}

	currentBuildsCount, err := s.planetStorage.GetCurrentBuildsCount(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetRepo.GetCurrentBuildsCount(): %w", err)
	}

	if currentBuildsCount >= consts.MaxBuildingsInProgress {
		return models.ErrTooManyBuildingsInProgress
	}

	err = s.recalcResources(ctx, planetID)
	if err != nil {
		return fmt.Errorf("recalcResources(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		var buildEvent models.BuildEvent

		currentBuildingID, err := planetRepo.GetBuildingID(ctx, planetID, consts.BuildingType(BuildingType))
		if err != nil {
			if !errors.Is(err, models.ErrBuildingNotFound) {
				return fmt.Errorf("planetRepo.GetBuildingID(): %w", err)
			}

			buildEvent, err = s.generateEventForNewBuilding(ctx, planetID, consts.BuildingType(BuildingType), planetRepo)
			if err != nil {
				return fmt.Errorf("generateEventForNewBuilding(): %w", err)
			}
		} else {
			buildEvent, err = s.generateEventForExistingBuilding(ctx, planetID, currentBuildingID, planetRepo)
		}

		err = planetRepo.SetFinishedBuildingTime(ctx, planetID, buildEvent.BuildingID, buildEvent.FinishedAt)
		if err != nil {
			return fmt.Errorf("planetRepo.SetFinishedBuildingTime(): %w", err)
		}

		err = planetRepo.CreateBuildingEvent(ctx, buildEvent)
		if err != nil {
			return fmt.Errorf("planetRepo.CreateBuildingEvent(): %w", err)
		}

		return nil
	})
}

func (s *Service) generateEventForNewBuilding(ctx context.Context, planetID uuid.UUID, buildingType consts.BuildingType, planetRepo TxStorages) (models.BuildEvent, error) {
	updatedAt := time.Now()

	stats, err := s.registry.GetBuildingZeroLvlStats(buildingType)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("registry.GetBuildingZeroLvlStats(): %w", err)
	}

	err = planetRepo.CreateBuilding(ctx, planetID, stats.ID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("planetRepo.CreateBuilding(): %w", err)
	}

	buildEvent := models.BuildEvent{
		PlanetID:   planetID,
		BuildingID: stats.ID,
		StartedAt:  updatedAt,
		FinishedAt: updatedAt.Add(time.Duration(stats.UpgradeTimeS) * time.Second),
	}

	return buildEvent, nil
}

func (s *Service) generateEventForExistingBuilding(ctx context.Context, planetID uuid.UUID, currentBuildingID consts.BuildingID, planetRepo TxStorages) (models.BuildEvent, error) {
	updatedAt := time.Now()

	// Calculate resources
	resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
	}

	updateBuildingStats, err := s.registry.GetBuildingNextLvlStats(currentBuildingID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("registry.GetBuildingStats(): %w", err)
	}

	if resources.Metal < updateBuildingStats.MetalCost ||
		resources.Crystal < updateBuildingStats.CrystalCost ||
		resources.Gas < updateBuildingStats.GasCost {
		return models.BuildEvent{}, models.ErrNotEnoughResources
	}

	leftResources := models.Resources{
		Metal:     resources.Metal - updateBuildingStats.MetalCost,
		Crystal:   resources.Crystal - updateBuildingStats.CrystalCost,
		Gas:       resources.Gas - updateBuildingStats.GasCost,
		UpdatedAt: updatedAt,
	}

	err = planetRepo.SetResources(ctx, planetID, leftResources)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("planetRepo.SetResources(): %w", err)
	}

	// get time to upgrade from registry
	currentStat, err := s.registry.GetBuildingStatsByID(currentBuildingID)
	if err != nil {
		return models.BuildEvent{}, fmt.Errorf("registry.GetBuildingStatsByID(): %w", err)
	}

	buildEvent := models.BuildEvent{
		PlanetID:   planetID,
		BuildingID: currentBuildingID,
		StartedAt:  updatedAt,
		FinishedAt: updatedAt.Add(time.Duration(currentStat.UpgradeTimeS) * time.Second),
	}

	return buildEvent, nil
}
