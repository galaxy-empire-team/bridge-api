package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Colonize(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, cargo models.Resources, fleet []models.PlanetFleetUnitCount) error {
	fleet = filterZeroCountFleet(fleet)

	if len(fleet) == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	if len(fleet) > s.registry.GetFleetUnitTypeCount() {
		return models.ErrInvalidInput
	}

	for _, fleetUnit := range fleet {
		if !s.registry.CheckFleetUnitIDExists(fleetUnit.ID) {
			return fmt.Errorf("%w: ID %d", models.ErrFleetIDNotExists, fleetUnit.ID)
		}
	}

	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, planetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if planetExists {
		return models.ErrColonizePlanetAlreadyExists
	}

	err = s.checkColonizationAvailability(ctx, userID)
	if err != nil {
		return fmt.Errorf("checkColonizationAvailability(): %w", err)
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeColonize)
	if err != nil {
		return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.updateResources(ctx, planetFrom, cargo, storages)
		if err != nil {
			return fmt.Errorf("updateResources(): %w", err)
		}

		err = s.updatePlanetFleet(ctx, planetFrom, fleet, storages)
		if err != nil {
			return fmt.Errorf("updatePlanetFleet(): %w", err)
		}

		var colonizeShipTaken bool
		for idx, fleetUnit := range fleet {
			// mission requires at least 1 colonizator that is removed here.
			stats, err := s.registry.GetFleetUnitStatsByID(fleetUnit.ID)
			if err != nil {
				return fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
			}

			if stats.Type == consts.FleetUnitTypeColonyShip && fleetUnit.Count > 0 {
				fleet[idx].Count -= 1
				colonizeShipTaken = true
			}
		}

		if !colonizeShipTaken {
			return models.ErrFleetNotFound
		}

		startedAt := time.Now().UTC()
		colonizeEvent := models.MissionEvent{
			UserID:      userID,
			PlanetFrom:  planetFrom,
			PlanetTo:    planetTo,
			Type:        missionID,
			Fleet:       fleet,
			Cargo:       cargo,
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

	researchIDs, err := s.researchStorage.GetUserResearches(ctx, userID)
	if err != nil {
		return fmt.Errorf("userStorage.GetUserResearches(): %w", err)
	}

	for _, researchID := range researchIDs {
		research, err := s.registry.GetResearchStatsByID(researchID)
		if err != nil {
			return fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
		}

		if research.Type != consts.ResearchTypeColonizeTechnology {
			continue
		}

		if research.Bonuses.AvaliableColonizePlanetCount == 0 {
			continue
		}

		if research.Bonuses.AvaliableColonizePlanetCount <= planetCount {
			return fmt.Errorf("planets count %d, availiable %d: %w", planetCount, research.Bonuses.AvaliableColonizePlanetCount, models.ErrColonizationNotAvailable)
		}

		return nil
	}

	return fmt.Errorf("%w", models.ErrColonizationNotAvailable)
}
