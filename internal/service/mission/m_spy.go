package mission

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Spy(ctx context.Context, mission models.MissionStart) (models.UserMission, error) {
	if err := s.validateSpyFleet(ctx, mission.Fleet, mission.Cargo); err != nil {
		return models.UserMission{}, fmt.Errorf("validateSpyFleet(): %w", err)
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

	spyMission, err := s.prepareUserMission(ctx, mission, consts.MissionTypeSpy)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("prepareUserMission(): %w", err)
	}

	return spyMission, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.removeFromPlanet(ctx, mission.UserID, mission.PlanetFrom, mission.Fleet, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("removeFromPlanet(): %w", err)
		}

		spyEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			MissionID:   spyMission.MissionID,
			Fleet:       mission.Fleet,
			IsReturning: spyMission.IsReturning,
			StartedAt:   spyMission.StartedAt,
			FinishedAt:  spyMission.FinishedAt,
		}

		eventID, err := storages.CreateMissionEvent(ctx, spyEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		spyMission.ID = eventID

		return nil
	})
}

func (s *Service) validateSpyFleet(ctx context.Context, fleet []models.FleetUnitCount, cargo models.Resources) error {
	if !cargo.IsEmpty() {
		return models.ErrCargoIsNotEmpty
	}

	if len(fleet) != 1 {
		return models.ErrInvalidInput
	}

	if fleet[0].Count == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	fleetUnitStats, err := s.registry.GetFleetUnitStatsByID(fleet[0].ID)
	if err != nil {
		return fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}
	if fleetUnitStats.Type != consts.FleetUnitTypeScout {
		return models.ErrInvalidShipTypeForSpyMission
	}

	return nil
}
