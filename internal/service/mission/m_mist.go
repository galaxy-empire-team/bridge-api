package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Mist(ctx context.Context, mission models.MissionStart) (models.UserMission, error) {
	if err := s.validateMistFleet(mission.Fleet, mission.Cargo); err != nil {
		return models.UserMission{}, fmt.Errorf("validateMistFleet(): %w", err)
	}

	if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, mission.PlanetFrom); err != nil {
		return models.UserMission{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	if mission.PlanetTo.Z != consts.MistPlanetCoordinateZ {
		return models.UserMission{}, models.ErrMistNotFound
	}

	mistMission, err := s.prepareUserMission(ctx, mission, consts.MissionTypeMist)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("prepareUserMission(): %w", err)
	}

	return mistMission, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err = s.removeFromPlanet(ctx, mission.UserID, mission.PlanetFrom, mission.Fleet, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("removeFromPlanet(): %w", err)
		}

		mistEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			MissionID:   mistMission.MissionID,
			Fleet:       mission.Fleet,
			Cargo:       mission.Cargo,
			IsReturning: mistMission.IsReturning,
			StartedAt:   mistMission.StartedAt,
			FinishedAt:  mistMission.FinishedAt,
		}

		eventID, err := storages.CreateMissionEvent(ctx, mistEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		mistMission.ID = eventID

		return nil
	})
}

func (s *Service) validateMistFleet(fleet []models.FleetUnitCount, cargo models.Resources) error {
	if !cargo.IsEmpty() {
		return models.ErrCargoIsNotEmpty
	}

	err := s.validateFleet(fleet, cargo)
	if err != nil {
		return fmt.Errorf("validateFleet(): %w", err)
	}

	return nil
}

func diffPositionY(x consts.PlanetPositionY, y consts.PlanetPositionY) consts.PlanetPositionY {
	if x < y {
		return y - x
	}

	return x - y
}
