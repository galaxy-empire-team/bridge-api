package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

const (
	// mistMissionDuration = 59 * time.Minute
	mistMissionDuration = 1 * time.Second
)

func (s *Service) Mist(ctx context.Context, mission models.MissionStart) error {
	if !mission.Cargo.IsEmpty() {
		return models.ErrCargoIsNotEmpty
	}

	if mission.PlanetTo.Z != consts.MistPlanetCoordinateZ {
		return models.ErrMistNotFound
	}

	err := s.checkFleetValid(mission.Fleet)
	if err != nil {
		return fmt.Errorf("checkFleetValid(): %w", err)
	}

	if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, mission.PlanetFrom); err != nil {
		return fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	planetCoordinates, err := s.planetStorage.GetCoordinates(ctx, mission.PlanetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.GetCoordinates(): %w", err)
	}

	// 1 is when a player start from current system
	distanceToMist := diffPositionY(planetCoordinates.Y, mission.PlanetTo.Y) + 1

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeMist)
	if err != nil {
		return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err = s.removeFleetFromPlanet(ctx, mission.PlanetFrom, mission.Fleet, storages)
		if err != nil {
			return fmt.Errorf("removeFleetFromPlanet(): %w", err)
		}

		startedAt := time.Now().UTC()
		mistEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			Type:        missionID,
			Fleet:       mission.Fleet,
			Cargo:       mission.Cargo,
			IsReturning: false,
			StartedAt:   startedAt,
			FinishedAt:  startedAt.Add(mistMissionDuration * time.Duration(distanceToMist)),
		}

		err = storages.CreateMissionEvent(ctx, mistEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		return nil
	})
}

func diffPositionY(x consts.PlanetPositionY, y consts.PlanetPositionY) consts.PlanetPositionY {
	if x < y {
		return y - x
	}

	return x - y
}
