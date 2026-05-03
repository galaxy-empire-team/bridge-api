package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Transport(ctx context.Context, mission models.MissionStart) error {
	err := s.checkFleetValid(mission.Fleet)
	if err != nil {
		return fmt.Errorf("checkFleetValid(): %w", err)
	}

	if !s.checkTransportCapacity(mission.Cargo, mission.Fleet, s.registry) {
		return models.ErrTransportCargoExceedsFleetCapacity
	}

	planetToID, err := s.planetStorage.GetIDByCoordinates(ctx, mission.PlanetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.GetIDByCoordinates(): %w", err)
	}

	for _, planetID := range []uuid.UUID{mission.PlanetFrom, planetToID} {
		isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, mission.UserID, planetID)
		if err != nil {
			return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
		}
		if !isUserPlanet {
			return models.ErrPlanetDoesNotBelongToUser
		}
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeTransport)
	if err != nil {
		return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	planetFromCoordinates, err := s.planetStorage.GetCoordinates(ctx, mission.PlanetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.GetCoordinates(): %w", err)
	}

	missionDuration, err := s.calculateMissionDuration(planetFromCoordinates, mission.PlanetTo, mission.Fleet, mission.SpeedMultiplier)
	if err != nil {
		return fmt.Errorf("calculateMissionDuration(): %w", err)
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.updateResources(ctx, mission.PlanetFrom, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("updateResources(): %w", err)
		}

		err = s.removeFleetFromPlanet(ctx, mission.PlanetFrom, mission.Fleet, storages)
		if err != nil {
			return fmt.Errorf("removeFleetFromPlanet(): %w", err)
		}

		startedAt := time.Now().UTC()
		transportEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			Type:        missionID,
			Fleet:       mission.Fleet,
			Cargo:       mission.Cargo,
			IsReturning: false,
			StartedAt:   startedAt,
			FinishedAt:  startedAt.Add(missionDuration),
		}

		err = storages.CreateMissionEvent(ctx, transportEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}
