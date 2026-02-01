package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
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
		planetBuildingsInfo, err := s.getMinesInfo(ctx, planetID)
		if err != nil {
			return fmt.Errorf("GetMinesInfo(): %w", err)
		}

		resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		millisecondsSinceLastUpdate := uint64(updatedAt.Sub(resources.UpdatedAt).Milliseconds())

		updatedResources := models.Resources{
			Metal:     resources.Metal + millisecondsSinceLastUpdate*planetBuildingsInfo[consts.BuildingTypeMetalMine].ProductionS/1000,
			Crystal:   resources.Crystal + millisecondsSinceLastUpdate*planetBuildingsInfo[consts.BuildingTypeCrystalMine].ProductionS/1000,
			Gas:       resources.Gas + millisecondsSinceLastUpdate*planetBuildingsInfo[consts.BuildingTypeGasMine].ProductionS/1000,
			UpdatedAt: updatedAt,
		}

		err = planetRepo.SetResources(ctx, planetID, updatedResources)
		if err != nil {
			return fmt.Errorf("planetRepo.SetResources(): %w", err)
		}

		return nil
	})
}

func (s *Service) getMinesInfo(ctx context.Context, planetID uuid.UUID) (map[consts.BuildingType]models.BuildingInfo, error) {
	mines, err := s.planetStorage.GetBuildingsInfo(ctx, planetID, consts.GetMineTypes())
	if err != nil {
		return nil, fmt.Errorf("planetRepo.GetBuildingsInfo(): %w", err)
	}

	// If mines are not build yet, initialize them with default values
	// TODO: get values from db
	for _, mineType := range consts.GetMineTypes() {
		if _, exists := mines[mineType]; !exists {
			mines[mineType] = models.BuildingInfo{
				Type:        mineType,
				Level:       defaultLvl,
				ProductionS: defaultProductionS,
			}
		}
	}

	return mines, nil
}
