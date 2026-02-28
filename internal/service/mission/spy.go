package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Spy(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, fleet []models.PlanetFleetUnitCount) error {
	fleet = filterZeroCountFleet(fleet)

	if len(fleet) == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	if len(fleet) != 1 {
		return models.ErrInvalidInput
	}

	fType, err := s.registry.GetFleetUnitStatsByID(fleet[0].ID)
	if err != nil {
		return fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}

	if fType.Type != consts.FleetUnitTypeScout {
		return models.ErrInvalidShipTypeForSpyMission
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

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeSpy)
	if err != nil {
		return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.updatePlanetFleet(ctx, planetFrom, fleet, storages)
		if err != nil {
			return fmt.Errorf("updatePlanetFleet(): %w", err)
		}

		startedAt := time.Now().UTC()
		spyEvent := models.MissionEvent{
			UserID:      userID,
			PlanetFrom:  planetFrom,
			PlanetTo:    planetTo,
			Type:        missionID,
			Fleet:       fleet,
			IsReturning: false,
			StartedAt:   startedAt,
			FinishedAt:  startedAt.Add(missionDuration),
		}

		err = storages.CreateMissionEvent(ctx, spyEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}
