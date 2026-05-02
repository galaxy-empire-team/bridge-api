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
	fleet := filterZeroCountFleet(mission.Fleet)

	if len(fleet) == 0 {
		return models.ErrFleetCannotBeEmpty
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

	if len(fleet) > s.registry.GetFleetUnitTypeCount() {
		return models.ErrInvalidInput
	}

	for _, fleetUnit := range fleet {
		if !s.registry.CheckFleetUnitIDExists(fleetUnit.ID) {
			return fmt.Errorf("%w: ID %d", models.ErrFleetIDNotExists, fleetUnit.ID)
		}
	}

	if !s.checkTransportCapacity(mission.Cargo, fleet, s.registry) {
		return models.ErrTransportCargoExceedsFleetCapacity
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeTransport)
	if err != nil {
		return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	planetFromCoordinates, err := s.planetStorage.GetCoordinates(ctx, mission.PlanetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.GetCoordinates(): %w", err)
	}

	missionDuration, err := s.calculateMissionDuration(planetFromCoordinates, mission.PlanetTo, fleet, mission.SpeedMultiplier)
	if err != nil {
		return fmt.Errorf("calculateMissionDuration(): %w", err)
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.updateResources(ctx, mission.PlanetFrom, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("updateResources(): %w", err)
		}

		err = s.updateFleet(ctx, mission.PlanetFrom, fleet, storages)
		if err != nil {
			return fmt.Errorf("updateFleet(): %w", err)
		}

		startedAt := time.Now().UTC()
		transportEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			Type:        missionID,
			Fleet:       fleet,
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

func (s *Service) checkTransportCapacity(cargo models.Resources, fleet []models.FleetUnitCount, registry registryProvider) bool {
	var cargoLimit uint64
	for _, fleetUnit := range fleet {
		fStats, err := registry.GetFleetUnitStatsByID(fleetUnit.ID)
		if err != nil {
			return false
		}

		cargoLimit += fStats.CargoCapacity * fleetUnit.Count
	}

	return cargo.Metal+cargo.Crystal+cargo.Gas <= cargoLimit
}
