package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) BoostFleetConstruction(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, boost models.UserBoost) (models.EventFinishTime, error) {
	var event models.EventFinishTime
	return event, s.txManager.ExecEventTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		fleetConstructionEvent, err := planetRepo.GetFleetConstructionEventForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetFleetConstructionEventForUpdate(): %w", err)
		}

		boostDuration, err := s.useBoost(ctx, userID, boost, planetRepo)
		if err != nil {
			return fmt.Errorf("useBoost(): %w", err)
		}

		fleetConstructionEvent.FinishedAt = fleetConstructionEvent.FinishedAt.Add(-boostDuration)
		err = planetRepo.SetFleetConstructionFinishTime(ctx, fleetConstructionEvent)
		if err != nil {
			return fmt.Errorf("planetStorage.SetFleetConstructionFinishTime(): %w", err)
		}

		event.StartedAt = fleetConstructionEvent.StartedAt
		event.FinishedAt = fleetConstructionEvent.FinishedAt

		return nil
	})
}
