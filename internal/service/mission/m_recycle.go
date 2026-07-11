package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Recycle(ctx context.Context, mission models.MissionStart) (models.UserMission, error) {
	if err := s.validateRecycleFleet(mission.Fleet, mission.Cargo); err != nil {
		return models.UserMission{}, fmt.Errorf("validateRecycleFleet(): %w", err)
	}

	if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, mission.PlanetFrom); err != nil {
		return models.UserMission{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, mission.PlanetTo)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if !planetExists && !s.isPlanetNPC(mission.PlanetTo.Z) {
		return models.UserMission{}, models.ErrPlanetNotFound
	}

	debris, err := s.planetStorage.GetDebris(ctx, mission.PlanetTo)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("planetStorage.GetDebris(): %w", err)
	}
	if debris.Metal == 0 && debris.Crystal == 0 {
		return models.UserMission{}, models.ErrNoDebrisFound
	}

	recycleMission, err := s.prepareUserMission(ctx, mission, consts.MissionTypeRecycle)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("prepareUserMission(): %w", err)
	}

	return recycleMission, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.removeFromPlanet(ctx, mission.UserID, mission.PlanetFrom, mission.Fleet, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("removeFromPlanet(): %w", err)
		}

		recycleEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			MissionID:   recycleMission.MissionID,
			Fleet:       mission.Fleet,
			IsReturning: recycleMission.IsReturning,
			StartedAt:   recycleMission.StartedAt,
			FinishedAt:  recycleMission.FinishedAt,
		}

		eventID, err := storages.CreateMissionEvent(ctx, recycleEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		recycleMission.ID = eventID

		return nil
	})
}

func (s *Service) validateRecycleFleet(fleet []models.FleetUnitCount, cargo models.Resources) error {
	if len(fleet) != 1 {
		return models.ErrInvalidInput
	}

	if fleet[0].Count == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	if !cargo.IsEmpty() {
		return models.ErrCargoIsNotEmpty
	}

	fleetUnitStats, err := s.registry.GetFleetUnitStatsByID(fleet[0].ID)
	if err != nil {
		return fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}
	if fleetUnitStats.Type != consts.FleetUnitTypeRecycler {
		return models.ErrInvalidShipTypeForRecycleMission
	}

	return nil
}
