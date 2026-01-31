package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

// recalcResources recalculates the resources of a planet based on the time since the last update.
// Recalcs using time.Now().UTC(). Use this before any operation that changes resources.
func (s *Service) recalcResources(ctx context.Context, planetID uuid.UUID) error {
	return s.recalcResourcesWithUpdatedTime(ctx, planetID, time.Now().UTC())
}

// recalcResourcesWithUpdatedTime recalculates the resources of a planet based on the time since the last update.
// Recalcs using the provided updatedAt time. Use this before any operation that changes resources.
func (s *Service) recalcResourcesWithUpdatedTime(ctx context.Context, planetID uuid.UUID, updatedAt time.Time) error {
	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		planetBuildingsInfo, err := planetRepo.GetBuildingsInfo(ctx, planetID, []models.BuildingType{
			models.BuildingTypeMetalMine,
			models.BuildingTypeCrystalMine,
			models.BuildingTypeGasMine,
		})
		if err != nil {
			return fmt.Errorf("planetRepo.GetBuildsInfo(): %w", err)
		}

		resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		millisecondsSinceLastUpdate := uint64(updatedAt.Sub(resources.UpdatedAt).Milliseconds())

		updatedResources := models.Resources{
			Metal:     resources.Metal + millisecondsSinceLastUpdate*planetBuildingsInfo[models.BuildingTypeMetalMine].ProductionS/1000,
			Crystal:   resources.Crystal + millisecondsSinceLastUpdate*planetBuildingsInfo[models.BuildingTypeCrystalMine].ProductionS/1000,
			Gas:       resources.Gas + millisecondsSinceLastUpdate*planetBuildingsInfo[models.BuildingTypeGasMine].ProductionS/1000,
			UpdatedAt: updatedAt,
		}

		err = planetRepo.SetResources(ctx, planetID, updatedResources)
		if err != nil {
			return fmt.Errorf("planetRepo.SetResources(): %w", err)
		}

		return nil
	})
}
