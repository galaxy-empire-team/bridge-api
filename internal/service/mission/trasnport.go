package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Transport(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, cargo models.Resources, fleet []models.PlanetFleetUnitCount) error {
	fleet = filterZeroCountFleet(fleet)

	if len(fleet) == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	planetToID, err := s.planetStorage.GetIDByCoordinates(ctx, planetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.GetIDByCoordinates(): %w", err)
	}

	for range []uuid.UUID{planetFrom, planetToID} {
		isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetFrom)
		if err != nil {
			return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
		}
		if !isUserPlanet {
			return models.ErrPlanetDoesNotBelongToUser
		}
	}

	// To prevent a way to attack I check the length of the fleet. I assume that client always sends the correct data.
	if len(fleet) > s.registry.GetFleetUnitTypeCount() {
		return models.ErrInvalidInput
	}

	for _, fleetUnit := range fleet {
		if !s.registry.CheckFleetUnitIDExists(fleetUnit.ID) {
			return fmt.Errorf("%w: ID %d", models.ErrFleetIDNotExists, fleetUnit.ID)
		}
	}

	if !s.checkTransportCapacity(cargo, fleet, s.registry) {
		return models.ErrTransportCargoExceedsFleetCapacity
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeTransport)
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

		startedAt := time.Now().UTC()
		spyEvent := models.MissionEvent{
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

		err = storages.CreateMissionEvent(ctx, spyEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}

func (s *Service) checkTransportCapacity(cargo models.Resources, fleet []models.PlanetFleetUnitCount, registry registryProvider) bool {
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
