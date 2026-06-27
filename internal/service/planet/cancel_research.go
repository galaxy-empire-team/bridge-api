package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) CancelResearch(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, researchID consts.ResearchID) error {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	return s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		err := planetRepo.DeleteResearchEvent(ctx, userID, researchID)
		if err != nil {
			return fmt.Errorf("planetStorage.DeleteBuildingEvent(): %w", err)
		}

		nextLvlResearchID, err := s.registry.GetResearchNextLvlID(researchID)
		if err != nil {
			return fmt.Errorf("registry.GetResearchNextLvlID(): %w", err)
		}

		nextLvlResearchStats, err := s.registry.GetResearchStatsByID(nextLvlResearchID)
		if err != nil {
			return fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
		}

		err = planetRepo.AddResources(ctx, planetID, models.Resources{
			Metal:   nextLvlResearchStats.MetalCost,
			Crystal: nextLvlResearchStats.CrystalCost,
			Gas:     nextLvlResearchStats.GasCost,
		})
		if err != nil {
			return fmt.Errorf("planetStorage.AddResources(): %w", err)
		}

		return nil
	})
}
