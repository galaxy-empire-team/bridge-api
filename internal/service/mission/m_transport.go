package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Transport(ctx context.Context, mission models.MissionStart) (models.UserMission, error) {
	err := s.validateFleet(mission.Fleet, mission.Cargo)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("checkFleetValid(): %w", err)
	}

	planetToID, err := s.planetStorage.GetIDByCoordinates(ctx, mission.PlanetTo)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("planetStorage.GetIDByCoordinates(): %w", err)
	}

	for _, planetID := range []uuid.UUID{mission.PlanetFrom, planetToID} {
		if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, planetID); err != nil {
			return models.UserMission{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
		}
	}

	transportMission, err := s.prepareUserMission(ctx, mission, consts.MissionTypeTransport)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("prepareUserMission(): %w", err)
	}

	return transportMission, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.removeFromPlanet(ctx, mission.UserID, mission.PlanetFrom, mission.Fleet, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("removeFromPlanet(): %w", err)
		}

		transportEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			MissionID:   transportMission.MissionID,
			Fleet:       mission.Fleet,
			Cargo:       mission.Cargo,
			IsReturning: transportMission.IsReturning,
			StartedAt:   transportMission.StartedAt,
			FinishedAt:  transportMission.FinishedAt,
		}

		eventID, err := storages.CreateMissionEvent(ctx, transportEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		transportMission.ID = eventID

		return nil
	})
}
