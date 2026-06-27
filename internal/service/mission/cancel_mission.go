package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) CancelMission(ctx context.Context, userID uuid.UUID, missionID uint64) (models.CancelMission, error) {
	updatedMissionEvent := models.CancelMission{}

	return updatedMissionEvent, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		missionEvent, err := storages.GetMissionForUpdate(ctx, userID, missionID)
		if err != nil {
			return fmt.Errorf("missionStorage.GetMissionForUpdate(): %w", err)
		}

		if missionEvent.IsReturning {
			return models.ErrMissionIsReturning
		}

		updatedMissionEvent = models.CancelMission{
			ID:          missionEvent.ID,
			IsReturning: true,
			StartedAt:   time.Now().UTC(),
			FinishedAt:  time.Now().UTC().Add(missionEvent.FinishedAt.Sub(missionEvent.StartedAt)),
		}

		err = storages.CancelMissionEvent(ctx, updatedMissionEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CancelMissionEvent(): %w", err)
		}

		return nil
	})
}
