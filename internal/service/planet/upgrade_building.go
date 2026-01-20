package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (s *Service) UpgradeBuilding(ctx context.Context, planetID uuid.UUID, buildType string) error {
	err := s.recalcResources(ctx, planetID)
	if err != nil {
		return fmt.Errorf("recalcResources(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		updatedTime := time.Now()
		planetBuild := models.PlanetBuildInfo{
			Type:      models.BuildType(buildType),
			Level:     0,
			UpdatedAt: updatedTime,
		}

		resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		currentBuildLvl, err := planetRepo.GetBuildLvl(ctx, planetID, planetBuild.Type)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		planetBuild.Level = currentBuildLvl

		// get stats for next level
		updatedBuildStats, err := planetRepo.GetBuildStats(ctx, planetBuild.Type, planetBuild.Level+1)
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildStats(): %w", err)
		}

		if resources.Metal <= updatedBuildStats.MetalCost ||
			resources.Crystal <= updatedBuildStats.CrystalCost ||
			resources.Gas <= updatedBuildStats.GasCost {
			return fmt.Errorf("not enough resources to build %s", buildType)
		}

		resources.Metal -= updatedBuildStats.MetalCost
		resources.Crystal -= updatedBuildStats.CrystalCost
		resources.Gas -= updatedBuildStats.GasCost
		resources.UpdatedAt = updatedTime

		err = planetRepo.SaveResources(ctx, planetID, resources)
		if err != nil {
			return fmt.Errorf("planetRepo.SaveResources(): %w", err)
		}

		planetBuild.Level++

		err = planetRepo.SaveBuild(ctx, planetID, planetBuild)
		if err != nil {
			return fmt.Errorf("planetRepo.SaveBuild(): %w", err)
		}

		return nil
	})
}
