package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Attack(ctx context.Context, userID uuid.UUID, planetFrom uuid.UUID, planetTo models.Coordinates, fleet []models.PlanetFleetUnitCount) error {
	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, planetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if !planetExists {
		return models.ErrAttackPlanetNotFound
	}

	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.ErrPlanetDoesNotBelongToUser
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

	// Make a map for the future before the transaction.
	reqFleet := lo.Associate(fleet, func(fleetUnit models.PlanetFleetUnitCount) (consts.FleetUnitID, uint64) {
		return fleetUnit.ID, fleetUnit.Count
	})

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		planetFleet, err := storages.GetFleetForUpdate(ctx, planetFrom)
		if err != nil {
			return fmt.Errorf("planetStorage.GetFleetCountForUpdate(): %w", err)
		}

		var leftFleetUnits []models.PlanetFleetUnitCount
		for _, fleetUnit := range planetFleet {
			reqCount, ok := reqFleet[fleetUnit.ID]
			if !ok {
				continue
			}

			if fleetUnit.Count < reqCount {
				return models.ErrNotEnoughFleetUnits
			}

			leftFleetUnits = append(leftFleetUnits, models.PlanetFleetUnitCount{
				ID:    fleetUnit.ID,
				Count: fleetUnit.Count - reqCount,
			})
		}

		err = storages.SetFleet(ctx, planetFrom, leftFleetUnits)
		if err != nil {
			return fmt.Errorf("planetStorage.SetFleet(): %w", err)
		}

		startedAt := time.Now().UTC()
		finishedAt := startedAt.Add(missionDuration)
		missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeAttack)
		if err != nil {
			return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
		}

		attackEvent := models.MissionEvent{
			UserID:      userID,
			PlanetFrom:  planetFrom,
			PlanetTo:    planetTo,
			Type:        missionID,
			Fleet:       fleet,
			IsReturning: false,
			StartedAt:   startedAt,
			FinishedAt:  finishedAt,
		}

		err = storages.CreateMissionEvent(ctx, attackEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}
