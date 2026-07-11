package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) CancelMission(ctx context.Context, userID uuid.UUID, missionID uint64) (models.CancelMission, error) {
	canceledMissionEvent := models.CancelMission{}

	return canceledMissionEvent, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		missionEvent, err := storages.GetMissionForUpdate(ctx, userID, missionID)
		if err != nil {
			return fmt.Errorf("missionStorage.GetMissionForUpdate(): %w", err)
		}

		if missionEvent.IsReturning {
			return models.ErrMissionIsReturning
		}

		planetFromCoordinates, err := s.planetStorage.GetCoordinates(ctx, missionEvent.PlanetFrom)
		if err != nil {
			return fmt.Errorf("planetStorage.GetCoordinates(): %w", err)
		}

		now := time.Now().UTC()

		canceledMissionEvent = models.CancelMission{
			ID:          missionEvent.ID,
			PlanetFrom:  uuid.Nil,
			PlanetTo:    planetFromCoordinates,
			IsReturning: true,
			StartedAt:   now,
			FinishedAt:  now.Add(now.Sub(missionEvent.StartedAt)),
		}

		err = storages.CancelMissionEvent(ctx, canceledMissionEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CancelMissionEvent(): %w", err)
		}

		return nil
	})
}
