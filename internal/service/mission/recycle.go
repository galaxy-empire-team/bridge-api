package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Recycle(ctx context.Context, mission models.MissionStart) error {
	if len(mission.Fleet) != 1 {
		return models.ErrInvalidInput
	}

	if mission.Fleet[0].Count == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	if mission.Cargo.IsEmpty() {
		return models.ErrCargoIsNotEmpty
	}

	fleetUnitStats, err := s.registry.GetFleetUnitStatsByID(mission.Fleet[0].ID)
	if err != nil {
		return fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}
	if fleetUnitStats.Type != consts.FleetUnitTypeRecycler {
		return models.ErrInvalidShipTypeForRecycleMission
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, mission.PlanetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if !planetExists && !s.isPlanetNPC(mission.PlanetTo.Z) {
		return models.ErrPlanetNotFound
	}

	if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, mission.PlanetFrom); err != nil {
		return fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	debris, err := s.planetStorage.GetDebris(ctx, mission.PlanetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.GetDebris(): %w", err)
	}
	if debris.Metal == 0 && debris.Crystal == 0 {
		return models.ErrNoDebrisFound
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeRecycle)
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
		err := s.removeFleetFromPlanet(ctx, mission.PlanetFrom, mission.Fleet, storages)
		if err != nil {
			return fmt.Errorf("removeFleetFromPlanet(): %w", err)
		}

		startedAt := time.Now().UTC()
		recycleEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			Type:        missionID,
			Fleet:       mission.Fleet,
			IsReturning: false,
			StartedAt:   startedAt,
			FinishedAt:  startedAt.Add(missionDuration*0 + 1*time.Second),
		}

		err = storages.CreateMissionEvent(ctx, recycleEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}
