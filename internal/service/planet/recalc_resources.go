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

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		planetBuildingsInfo, err := s.getMinesInfo(ctx, planetID)
		if err != nil {
			return fmt.Errorf("GetMinesInfo(): %w", err)
		}

		resources, err := planetRepo.GetResourcesForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetRepo.GetResourcesForUpdate(): %w", err)
		}

		millisecondsSinceLastUpdate := updatedAt.Sub(resources.UpdatedAt).Milliseconds()
		if millisecondsSinceLastUpdate <= 0 {
			return nil
		}

		metalProductionPerSecond := float32(planetBuildingsInfo[consts.BuildingTypeMetalMine].ProductionS) * multiplier
		crystalProductionPerSecond := float32(planetBuildingsInfo[consts.BuildingTypeCrystalMine].ProductionS) * multiplier
		gasProductionPerSecond := float32(planetBuildingsInfo[consts.BuildingTypeGasMine].ProductionS) * multiplier

		updatedResources := models.Resources{
			Metal:     resources.Metal + uint64(millisecondsSinceLastUpdate)*uint64(metalProductionPerSecond)/1000,
			Crystal:   resources.Crystal + uint64(millisecondsSinceLastUpdate)*uint64(crystalProductionPerSecond)/1000,
			Gas:       resources.Gas + uint64(millisecondsSinceLastUpdate)*uint64(gasProductionPerSecond)/1000,
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
	for _, mineType := range consts.GetMineTypes() {
		if _, exists := mines[mineType]; !exists {
			stat, err := s.registry.GetBuildingZeroLvlStats(mineType)
			if err != nil {
				return nil, fmt.Errorf("registry.GetBuildingZeroLvlStats(): %w", err)
			}

			mines[mineType] = models.BuildingInfo{
				Type:        stat.Type,
				Level:       stat.Level,
				ProductionS: stat.ProductionS,
			}
		}
	}

	return mines, nil
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
