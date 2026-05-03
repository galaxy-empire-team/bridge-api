package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Colonize(ctx context.Context, mission models.MissionStart) error {
	if len(mission.Fleet) != 1 {
		return models.ErrInvalidInput
	}

	if mission.Fleet[0].Count == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	fleetUnitStats, err := s.registry.GetFleetUnitStatsByID(mission.Fleet[0].ID)
	if err != nil {
		return fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}
	if fleetUnitStats.Type != consts.FleetUnitTypeColonyShip {
		return models.ErrInvalidShipTypeForColonization
	}

	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, mission.UserID, mission.PlanetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, mission.PlanetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if planetExists {
		return models.ErrColonizePlanetAlreadyExists
	}

	err = s.checkColonizationAvailability(ctx, mission.UserID)
	if err != nil {
		return fmt.Errorf("checkColonizationAvailability(): %w", err)
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeColonize)
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

		var colonizeShipTaken bool
		for idx, fleetUnit := range mission.Fleet {
			// mission requires at least 1 colonizator that is removed here.
			stats, err := s.registry.GetFleetUnitStatsByID(fleetUnit.ID)
			if err != nil {
				return fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
			}

			if stats.Type == consts.FleetUnitTypeColonyShip && fleetUnit.Count > 0 {
				mission.Fleet[idx].Count -= 1
				colonizeShipTaken = true
			}
		}

		if !colonizeShipTaken {
			return models.ErrFleetNotFound
		}

		startedAt := time.Now().UTC()
		colonizeEvent := models.MissionEvent{
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

		err = storages.CreateMissionEvent(ctx, colonizeEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}

func (s *Service) checkColonizationAvailability(ctx context.Context, userID uuid.UUID) error {
	planetCount, err := s.planetStorage.GetUserPlanetsCount(ctx, userID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetUserPlanetsCount(): %w", err)
	}

	userResearches, err := s.researchStorage.GetUserResearchesByTypes(ctx, userID, []consts.ResearchType{consts.ResearchTypeColonizeTechnology})
	if err != nil {
		return fmt.Errorf("userStorage.GetUserResearches(): %w", err)
	}

	researchID, ok := userResearches[consts.ResearchTypeColonizeTechnology]
	if !ok {
		return models.ErrColonizationNotAvailable
	}

	researchStats, err := s.registry.GetResearchStatsByID(researchID)
	if err != nil {
		return fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
	}

	if planetCount >= researchStats.Bonuses.AvaliableColonizePlanetCount {
		return models.ErrColonizationNotAvailable
	}

	return nil
}
