package planet

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
func (s *Repository) RecalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error {
	return s.RecalcResourcesWithUpdatedTime(ctx, userID, planetID, time.Now().UTC())
}

// RecalcResourcesWithUpdatedTime recalculates the resources of a planet based on the time since the last update.
// Recalcs using the provided updatedAt time. Use this before any operation that changes resources
func (s *Repository) RecalcResourcesWithUpdatedTime(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedAt time.Time) error {
	multiplier, err := s.getResearchResourceMultiplier(ctx, userID)
	if err != nil {
		return fmt.Errorf("getResearchResourceMultiplier(): %w", err)
	}

	mines, err := s.planetStorage.GetPlanetMinesProduction(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetPlanetMinesProduction(): %w", err)
	}

	return s.txManager.ExecResourceRepoTx(ctx, func(ctx context.Context, storage TxStorages) error {
		resources, err := storage.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("storage.GetResourcesForUpdate(): %w", err)
		}

		millisecondsSinceLastUpdate := updatedAt.Sub(resources.UpdatedAt).Milliseconds()
		if millisecondsSinceLastUpdate <= 0 {
			return nil
		}

		metalProductionPerSecond := float32(mines[consts.BuildingTypeMetalMine]) * multiplier
		crystalProductionPerSecond := float32(mines[consts.BuildingTypeCrystalMine]) * multiplier
		gasProductionPerSecond := float32(mines[consts.BuildingTypeGasMine]) * multiplier

		updatedResources := models.Resources{
			Metal:     resources.Metal + uint64(millisecondsSinceLastUpdate)*uint64(metalProductionPerSecond)/1000,
			Crystal:   resources.Crystal + uint64(millisecondsSinceLastUpdate)*uint64(crystalProductionPerSecond)/1000,
			Gas:       resources.Gas + uint64(millisecondsSinceLastUpdate)*uint64(gasProductionPerSecond)/1000,
			UpdatedAt: updatedAt,
		}

		err = storage.SetResources(ctx, planetID, updatedResources)
		if err != nil {
			return fmt.Errorf("storage.SetResources(): %w", err)
		}

		return nil
	})
}
