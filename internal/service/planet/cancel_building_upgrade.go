package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) CancelBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) error {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		err := planetRepo.DeleteBuildingEvent(ctx, planetID, buildingID)
		if err != nil {
			return fmt.Errorf("planetStorage.DeleteBuildingEvent(): %w", err)
		}

		nextLvlBuildingID, err := s.registry.GetBuildingNextLvlID(buildingID)
		if err != nil {
			return fmt.Errorf("registry.GetBuildingNextLvlID(): %w", err)
		}

		nextLvlStats, err := s.registry.GetBuildingStatsByID(nextLvlBuildingID)
		if err != nil {
			return fmt.Errorf("registry.GetBuildingStatsByID(): %w", err)
		}

		err = planetRepo.AddResources(ctx, planetID, models.Resources{
			Metal:   nextLvlStats.MetalCost,
			Crystal: nextLvlStats.CrystalCost,
			Gas:     nextLvlStats.GasCost,
		})
		if err != nil {
			return fmt.Errorf("planetStorage.AddResources(): %w", err)
		}

		return nil
	})
}
