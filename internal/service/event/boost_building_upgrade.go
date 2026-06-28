package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) BoostBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID, boost models.UserBoost) (models.EventFinishTime, error) {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return models.EventFinishTime{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	var event models.EventFinishTime
	return event, s.txManager.ExecEventTx(ctx, func(ctx context.Context, planetRepo TxStorages) error {
		buildingEvent, err := planetRepo.GetBuildingEventForUpdate(ctx, planetID, buildingID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetBuildingEventForUpdate(): %w", err)
		}

		boostDuration, err := s.useBoost(ctx, userID, boost, planetRepo)
		if err != nil {
			return fmt.Errorf("useBoost(): %w", err)
		}

		buildingEvent.FinishedAt = buildingEvent.FinishedAt.Add(-boostDuration)
		err = planetRepo.SetBuildingFinishTime(ctx, buildingEvent)
		if err != nil {
			return fmt.Errorf("planetStorage.SetBuildingFinishTime(): %w", err)
		}

		event.StartedAt = buildingEvent.StartedAt
		event.FinishedAt = buildingEvent.FinishedAt

		return nil
	})
}
