package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

// recalcResources recalculates the resources of a planet based on the time since the last update
// use before any operation that changes resources.
func (s *Service) recalcResources(ctx context.Context, planetID uuid.UUID) error {
	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		planetBuildsInfo, err := planetRepo.GetBuildsInfo(ctx, planetID, []models.BuildType{
			models.BuildingTypeMetalMine,
			models.BuildingTypeCrystalMine,
			models.BuildingTypeGasMine,
		})
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildsInfo(): %w", err)
		}

		err = getDefaultIfMinesNotExists(ctx, planetBuildsInfo, planetRepo)
		if err != nil {
			return fmt.Errorf("getDefaultIfMinesNotExists(): %w", err)
		}

		resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		updatedTime := time.Now()
		secondsSinceLastUpdate := updatedTime.Sub(resources.UpdatedAt).Milliseconds()

		updatedResources := models.Resources{
			Metal:     resources.Metal + uint64(secondsSinceLastUpdate)*planetBuildsInfo[models.BuildingTypeMetalMine].MetalPerSecond/1000,
			Crystal:   resources.Crystal + uint64(secondsSinceLastUpdate)*planetBuildsInfo[models.BuildingTypeCrystalMine].CrystalPerSecond/1000,
			Gas:       resources.Gas + uint64(secondsSinceLastUpdate)*planetBuildsInfo[models.BuildingTypeGasMine].GasPerSecond/1000,
			UpdatedAt: updatedTime,
		}

		err = planetRepo.SaveResources(ctx, planetID, updatedResources)
		if err != nil {
			return fmt.Errorf("planetRepo.SaveResources(): %w", err)
		}

		return nil
	})
}

// getDefaultIfMinesNotExists adds default level 0 mines to the map if they do not exist
func getDefaultIfMinesNotExists(
	ctx context.Context,
	planetMines map[models.BuildType]models.PlanetBuildInfo,
	planetRepo TxStorages,
) error {
	buildTypes := []models.BuildType{
		models.BuildingTypeMetalMine,
		models.BuildingTypeCrystalMine,
		models.BuildingTypeGasMine,
	}

	for _, buildType := range buildTypes {
		if _, exists := planetMines[buildType]; exists {
			continue
		}

		mine, err := planetRepo.GetBuildStats(ctx, buildType, defaultLvl)
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildStats(): %w", err)
		}
		planetMines[buildType] = models.PlanetBuildInfo{
			Type:             mine.Type,
			Level:            mine.Level,
			MetalPerSecond:   mine.MetalPerSecond,
			CrystalPerSecond: mine.CrystalPerSecond,
			GasPerSecond:     mine.GasPerSecond,
		}
	}

	return nil
}
