package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) UpgradeBuilding(ctx context.Context, planetID uuid.UUID, BuildingType string) error {
	if !models.IsValidBuildingType(models.BuildingType(BuildingType)) {
		return models.ErrBuildTypeInvalid
	}

	currentBuildsCount, err := s.planetStorage.GetCurrentintBuildsCount(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetRepo.GetCurrentintBuildsCount(): %w", err)
	}

	if currentBuildsCount >= maxBuildingsInProgress {
		return models.ErrTooManyBuildingsInProgress
	}

	err = s.recalcResources(ctx, planetID)
	if err != nil {
		return fmt.Errorf("recalcResources(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		planetBuild := models.BuildingInfo{
			Type:      models.BuildingType(BuildingType),
			UpdatedAt: time.Now().UTC(),
		}

		resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		currentBuildLvl, err := planetRepo.GetBuildingLvl(ctx, planetID, planetBuild.Type)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		planetBuild.Level = currentBuildLvl

		if planetBuild.Level >= maxBuildingLvl {
			return models.ErrBuildingMaxLevelReached
		}

		// get stats for the next level
		updateBuildingStats, err := planetRepo.GetBuildingStats(ctx, planetBuild.Type, planetBuild.Level+1)
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildingStats(): %w", err)
		}

		if resources.Metal <= updateBuildingStats.MetalCost ||
			resources.Crystal <= updateBuildingStats.CrystalCost ||
			resources.Gas <= updateBuildingStats.GasCost {
			return models.ErrNotEngoughResources
		}

		resources.Metal -= updateBuildingStats.MetalCost
		resources.Crystal -= updateBuildingStats.CrystalCost
		resources.Gas -= updateBuildingStats.GasCost
		resources.UpdatedAt = planetBuild.UpdatedAt

		err = planetRepo.SetResources(ctx, planetID, resources)
		if err != nil {
			return fmt.Errorf("planetRepo.SetResources(): %w", err)
		}

		finishedAt := planetBuild.UpdatedAt.Add(time.Duration(updateBuildingStats.UpgradeTimeS) * time.Second)
		planetBuild.FinishedAt = finishedAt

		err = planetRepo.SetFinishedBuildingTime(ctx, planetID, planetBuild)
		if err != nil {
			return fmt.Errorf("planetRepo.SetFinishedBuildingTime(): %w", err)
		}

		buildingEvent := models.BuildEvent{
			PlanetID:     planetID,
			BuildingType: planetBuild.Type,
			FinishedAt:   finishedAt,
		}

		err = planetRepo.CreateBuildingEvent(ctx, buildingEvent)
		if err != nil {
			return fmt.Errorf("planetRepo.CreateBuildingEvent(): %w", err)
		}

		return nil
	})
}
