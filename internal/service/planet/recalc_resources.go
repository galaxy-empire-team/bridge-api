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
func (s *Service) recalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error {
	return s.recalcResourcesWithUpdatedTime(ctx, userID, planetID, time.Now().UTC())
}

// recalcResourcesWithUpdatedTime recalculates the resources of a planet based on the time since the last update.
// Recalcs using the provided updatedAt time. Use this before any operation that changes resources.
func (s *Service) recalcResourcesWithUpdatedTime(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedAt time.Time) error {
	multiplier, err := s.getResearchResourceMultiplier(ctx, userID)
	if err != nil {
		return fmt.Errorf("getResearchResourceMultiplier(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetStorage TxStorages) error {
		mines, err := s.planetStorage.GetPlanetMinesProduction(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetPlanetMinesProduction(): %w", err)
		}

		resources, err := s.planetStorage.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetResourcesForUpdate(): %w", err)
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

		err = planetStorage.SetResources(ctx, planetID, updatedResources)
		if err != nil {
			return fmt.Errorf("planetStorage.SetResources(): %w", err)
		}

		return nil
	})
}

func (s *Service) getResearchResourceMultiplier(ctx context.Context, userID uuid.UUID) (float32, error) {
	researchIDs, err := s.researchStorage.GetUserResearches(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("researchStorage.GetUserResearches(): %w", err)
	}

	for _, researchID := range researchIDs {
		research, err := s.registry.GetResearchStatsByID(researchID)
		if err != nil {
			return 0, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
		}

		if research.Type != consts.ResearchTypeIndustrialTechnology {
			continue
		}

		return research.Bonuses.ProductionSpeedImprove, nil
	}

	// If user has no industrial technology research, return 1
	return 1, nil
}
