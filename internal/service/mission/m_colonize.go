package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Colonize(ctx context.Context, mission models.MissionStart) (models.UserMission, error) {
	if err := s.validateColonizeFleet(mission.Fleet, mission.Cargo); err != nil {
		return models.UserMission{}, fmt.Errorf("validateColonizeFleet(): %w", err)
	}

	if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, mission.PlanetFrom); err != nil {
		return models.UserMission{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	err := s.checkColonizationAvailability(ctx, mission.UserID)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("checkColonizationAvailability(): %w", err)
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, mission.PlanetTo)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if planetExists {
		return models.UserMission{}, models.ErrColonizePlanetAlreadyExists
	}

	colonizeMission, err := s.prepareUserMission(ctx, mission, consts.MissionTypeColonize)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("prepareUserMission(): %w", err)
	}

	return colonizeMission, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err = s.removeFromPlanet(ctx, mission.UserID, mission.PlanetFrom, mission.Fleet, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("removeFromPlanet(): %w", err)
		}

		// remove one colonizator
		mission.Fleet[0].Count -= 1

		colonizeEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			MissionID:   colonizeMission.MissionID,
			Fleet:       mission.Fleet,
			Cargo:       mission.Cargo,
			IsReturning: colonizeMission.IsReturning,
			StartedAt:   colonizeMission.StartedAt,
			FinishedAt:  colonizeMission.FinishedAt,
		}

		eventID, err := storages.CreateMissionEvent(ctx, colonizeEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		colonizeMission.ID = eventID

		return nil
	})
}

func (s *Service) checkColonizationAvailability(ctx context.Context, userID uuid.UUID) error {
	planetCount, err := s.planetStorage.GetUserPlanetsCount(ctx, userID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetUserPlanetsCount(): %w", err)
	}

	colonizeResearchStats, err := s.repository.GetResearchByType(ctx, userID, consts.ResearchTypeColonizeTechnology)
	if err != nil {
		return fmt.Errorf("repository.GetResearchByType(): %w", err)
	}

	if planetCount >= colonizeResearchStats.Bonuses.AvaliableColonizePlanetCount {
		return models.ErrColonizationNotAvailable
	}

	return nil
}

func (s *Service) validateColonizeFleet(fleet []models.FleetUnitCount, cargo models.Resources) error {
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
	if fleetUnitStats.Type != consts.FleetUnitTypeColonyShip {
		return models.ErrInvalidShipTypeForColonization
	}

	return nil
}
