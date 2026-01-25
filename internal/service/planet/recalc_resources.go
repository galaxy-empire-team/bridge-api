package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

// recalcResources recalculates the resources of a planet based on the time since the last update
// use before any operation that changes resources.
func (s *Service) recalcResources(ctx context.Context, planetID uuid.UUID) error {
	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		planetBuildingsInfo, err := planetRepo.GetBuildingsInfo(ctx, planetID, []models.BuildingType{
			models.BuildingTypeMetalMine,
			models.BuildingTypeCrystalMine,
			models.BuildingTypeGasMine,
		})
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildsInfo(): %w", err)
		}

		err = getDefaultIfMinesNotExists(ctx, planetBuildingsInfo, planetRepo)
		if err != nil {
			return fmt.Errorf("getDefaultIfMinesNotExists(): %w", err)
		}

		resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		updatedTime := time.Now()
		secondsSinceLastUpdate := uint64(updatedTime.Sub(resources.UpdatedAt).Milliseconds())

		updatedResources := models.Resources{
			Metal:     resources.Metal + secondsSinceLastUpdate*planetBuildingsInfo[models.BuildingTypeMetalMine].MetalPerSecond/1000,
			Crystal:   resources.Crystal + secondsSinceLastUpdate*planetBuildingsInfo[models.BuildingTypeCrystalMine].CrystalPerSecond/1000,
			Gas:       resources.Gas + secondsSinceLastUpdate*planetBuildingsInfo[models.BuildingTypeGasMine].GasPerSecond/1000,
			UpdatedAt: updatedTime,
		}

		err = planetRepo.SetResources(ctx, planetID, updatedResources)
		if err != nil {
			return fmt.Errorf("planetRepo.SetResources(): %w", err)
		}

		return nil
	})
}

// getDefaultIfMinesNotExists adds default level 0 mines to the map if they do not exist
func getDefaultIfMinesNotExists(
	ctx context.Context,
	planetMines map[models.BuildingType]models.BuildingInfo,
	planetRepo TxStorages,
) error {
	BuildingTypes := []models.BuildingType{
		models.BuildingTypeMetalMine,
		models.BuildingTypeCrystalMine,
		models.BuildingTypeGasMine,
	}

	for _, BuildingType := range BuildingTypes {
		if _, exists := planetMines[BuildingType]; exists {
			continue
		}

		mine, err := planetRepo.GetBuildingStats(ctx, BuildingType, defaultLvl)
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildingStats(): %w", err)
		}
		planetMines[BuildingType] = models.BuildingInfo{
			Type:             mine.Type,
			Level:            mine.Level,
			MetalPerSecond:   mine.MetalPerSecond,
			CrystalPerSecond: mine.CrystalPerSecond,
			GasPerSecond:     mine.GasPerSecond,
		}
	}

	return nil
}
