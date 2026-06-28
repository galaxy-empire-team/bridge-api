package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) BoostResearch(ctx context.Context, userID uuid.UUID, researchID consts.ResearchID, boost models.UserBoost) (models.EventFinishTime, error) {
	var event models.EventFinishTime
	return event, s.txManager.ExecEventTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		researchEvent, err := planetRepo.GetResearchEventForUpdate(ctx, userID, researchID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetResearchEventForUpdate(): %w", err)
		}

		boostDuration, err := s.useBoost(ctx, userID, boost, planetRepo)
		if err != nil {
			return fmt.Errorf("useBoost(): %w", err)
		}

		researchEvent.FinishedAt = researchEvent.FinishedAt.Add(-boostDuration)
		err = planetRepo.SetResearchFinishTime(ctx, researchEvent)
		if err != nil {
			return fmt.Errorf("planetStorage.SetResearchFinishTime(): %w", err)
		}

		event.StartedAt = researchEvent.StartedAt
		event.FinishedAt = researchEvent.FinishedAt

		return nil
	})
}
