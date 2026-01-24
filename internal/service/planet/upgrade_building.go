package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (s *Service) UpgradeBuilding(ctx context.Context, planetID uuid.UUID, BuildingType string) error {
	err := s.recalcResources(ctx, planetID)
	if err != nil {
		return fmt.Errorf("recalcResources(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		updatedTime := time.Now().UTC()

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

		// get stats for the next level
		updateBuildingStats, err := planetRepo.GetBuildingStats(ctx, planetBuild.Type, planetBuild.Level+1)
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildingStats(): %w", err)
		}

		if resources.Metal <= updateBuildingStats.MetalCost ||
			resources.Crystal <= updateBuildingStats.CrystalCost ||
			resources.Gas <= updateBuildingStats.GasCost {
			return fmt.Errorf("not enough resources to build %s", BuildingType)
		}

		resources.Metal -= updateBuildingStats.MetalCost
		resources.Crystal -= updateBuildingStats.CrystalCost
		resources.Gas -= updateBuildingStats.GasCost
		resources.UpdatedAt = updatedTime

		err = planetRepo.SetResources(ctx, planetID, resources)
		if err != nil {
			return fmt.Errorf("planetRepo.SetResources(): %w", err)
		}

		finishedAt := updatedTime.Add(time.Duration(updateBuildingStats.UpgradeTimeInSeconds) * time.Second)
		planetBuild.FinishedAt = &finishedAt

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
