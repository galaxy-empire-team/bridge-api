package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Attack(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, fleet []models.PlanetFleetUnitCount) error {
	fleet = filterZeroCountFleet(fleet)

	if len(fleet) == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, planetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if !planetExists {
		return models.ErrPlanetNotFound
	}

	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
	}

	if len(fleet) > s.registry.GetFleetUnitTypeCount() {
		return models.ErrInvalidInput
	}

	for _, fleetUnit := range fleet {
		if !s.registry.CheckFleetUnitIDExists(fleetUnit.ID) {
			return fmt.Errorf("%w: ID %d", models.ErrFleetIDNotExists, fleetUnit.ID)
		}
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeAttack)
	if err != nil {
		return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err = s.updatePlanetFleet(ctx, planetFrom, fleet, storages)
		if err != nil {
			return fmt.Errorf("updatePlanetFleet(): %w", err)
		}

		startedAt := time.Now().UTC()
		attackEvent := models.MissionEvent{
			UserID:      userID,
			PlanetFrom:  planetFrom,
			PlanetTo:    planetTo,
			Type:        missionID,
			Fleet:       fleet,
			IsReturning: false,
			StartedAt:   startedAt,
			FinishedAt:  startedAt.Add(missionDuration),
		}

		err = storages.CreateMissionEvent(ctx, attackEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}
