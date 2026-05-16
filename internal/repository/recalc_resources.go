package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

// RecalcResources recalculates the resources of a planet based on the time since the last update.
// Recalcs using time.Now().UTC(). Use this before any operation that changes resources.
func (r *Repository) RecalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error {
	return r.RecalcResourcesWithUpdatedTime(ctx, userID, planetID, time.Now().UTC())
}

// RecalcResourcesWithUpdatedTime recalculates the resources of a planet based on the time since the last update.
// Recalcs using the provided updatedAt time. Use this before any operation that changes resources
func (r *Repository) RecalcResourcesWithUpdatedTime(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedAt time.Time) error {
	industrialTechStat, err := r.GetResearchByType(ctx, userID, consts.ResearchTypeIndustrialTechnology)
	if err != nil {
		return fmt.Errorf("r.GetResearchByType(): %w", err)
	}

	minesProduction, err := r.planetStorage.GetPlanetMinesProduction(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetPlanetMinesProduction(): %w", err)
	}

	// If building is not build yet I fill it with zero level production speed from registry
	const mineTypesCount = 3
	if len(minesProduction) != mineTypesCount {
		for _, buildingType := range []consts.BuildingType{consts.BuildingTypeMetalMine, consts.BuildingTypeCrystalMine, consts.BuildingTypeGasMine} {
			if _, ok := minesProduction[buildingType]; ok {
				continue
			}

			id, err := r.registry.GetBuildingZeroLvlIDByType(buildingType)
			if err != nil {
				return fmt.Errorf("registry.GetBuildingZeroLvlIDByType(): %w", err)
			}

			bStat, err := r.registry.GetBuildingStatsByID(id)
			if err != nil {
				return fmt.Errorf("registry.GetBuildingStatsByID(): %w", err)
			}

			minesProduction[buildingType] = bStat.ProductionS
		}
	}

	return r.txManager.ExecResourceRepoTx(ctx, func(ctx context.Context, storage TxStorages) error {
		resources, err := storage.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("storage.GetResourcesForUpdate(): %w", err)
		}

		millisecondsSinceLastUpdate := updatedAt.Sub(resources.UpdatedAt).Milliseconds()
		if millisecondsSinceLastUpdate <= 0 {
			return nil
		}

		metalProductionPerSecond := float32(minesProduction[consts.BuildingTypeMetalMine]) * industrialTechStat.Bonuses.ProductionSpeedImprove
		crystalProductionPerSecond := float32(minesProduction[consts.BuildingTypeCrystalMine]) * industrialTechStat.Bonuses.ProductionSpeedImprove
		gasProductionPerSecond := float32(minesProduction[consts.BuildingTypeGasMine]) * industrialTechStat.Bonuses.ProductionSpeedImprove

		metalIncrease := float64(millisecondsSinceLastUpdate) * float64(metalProductionPerSecond) / 1000
		crystalIncrease := float64(millisecondsSinceLastUpdate) * float64(crystalProductionPerSecond) / 1000
		gasIncrease := float64(millisecondsSinceLastUpdate) * float64(gasProductionPerSecond) / 1000

		updatedResources := models.Resources{
			Metal:     resources.Metal + uint64(metalIncrease),
			Crystal:   resources.Crystal + uint64(crystalIncrease),
			Gas:       resources.Gas + uint64(gasIncrease),
			UpdatedAt: updatedAt,
		}

		err = storage.SetResources(ctx, planetID, updatedResources)
		if err != nil {
			return fmt.Errorf("storage.SetResources(): %w", err)
		}

		return nil
	})
}
