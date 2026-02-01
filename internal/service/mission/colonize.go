package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Colonize(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates) error {
	// TODO after fleet implementation, check if user has any colonization fleets available
	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, planetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if planetExists {
		return models.ErrColonizePlanetAlreadyExists
	}

	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		startedAt := time.Now().UTC()
		finishedAt := startedAt
		colonizeEvent := models.MissionEvent{
			UserID:     userID,
			PlanetFrom: planetFrom,
			PlanetTo:   planetTo,
			Type:       consts.MissionTypeColonize,
			StartedAt:  startedAt,
			FinishedAt: finishedAt,
		}

		err = storages.CreateMissionEvent(ctx, colonizeEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}
