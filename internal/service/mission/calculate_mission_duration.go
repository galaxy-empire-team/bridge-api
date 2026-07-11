package mission

import (
	"fmt"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

const (
	// Addition time to prevent fast missions
	missionDurationAddSpy  = 15 * time.Second
	missionDurationAdd     = 3 * time.Minute
	missionDurationAddMist = 40 * time.Minute
)

type coordinate interface {
	consts.PlanetPositionX | consts.PlanetPositionY | consts.PlanetPositionZ
}

type durationInput struct {
	PlanetFrom      models.Coordinates
	PlanetTo        models.Coordinates
	Fleet           []models.FleetUnitCount
	SpeedMultiplier float64
	Type            consts.MissionType
}

func (s *Service) calculateMissionDuration(input durationInput) (time.Duration, error) {
	if input.SpeedMultiplier < consts.SpeedMultiplierMin || input.SpeedMultiplier > consts.SpeedMultiplierMax {
		return 0, fmt.Errorf("invalid speed multiplier: %f", input.SpeedMultiplier)
	}

	minSpeed, err := s.calcMinSpeed(input.Fleet)
	if err != nil {
		return 0, fmt.Errorf("calcMinSpeed(): %w", err)
	}

	duration := calcDistanceDuration(input.PlanetFrom, input.PlanetTo) / time.Duration(minSpeed)
	duration = time.Duration(float64(duration) / input.SpeedMultiplier)

	switch input.Type {
	case consts.MissionTypeSpy:
		duration += missionDurationAddSpy
	case consts.MissionTypeColonize:
		duration += missionDurationAdd
	case consts.MissionTypeMist:
		duration += missionDurationAddMist
	}

	return duration, nil
}

func (s *Service) calcMinSpeed(fleet []models.FleetUnitCount) (uint64, error) {
	if len(fleet) == 0 {
		return 0, models.ErrFleetCannotBeEmpty
	}

	unitStats, err := s.registry.GetFleetUnitStatsByID(fleet[0].ID)
	if err != nil {
		return 0, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
	}

	minSpeed := unitStats.Speed

	for _, unit := range fleet[1:] {
		unitStats, err := s.registry.GetFleetUnitStatsByID(unit.ID)
		if err != nil {
			return 0, fmt.Errorf("registry.GetFleetUnitStatsByID(): %w", err)
		}

		if unitStats.Speed < minSpeed {
			minSpeed = unitStats.Speed
		}
	}

	return minSpeed, nil
}

func calcDistanceDuration(from models.Coordinates, to models.Coordinates) time.Duration {
	var timeToReachPlanet uint64

	timeToReachPlanet += consts.TimeToReachNearestGalaxyS * uint64(cooridanteDiff(from.X, to.X))
	timeToReachPlanet += consts.TimeToReachNearestSystemS * uint64(cooridanteDiff(from.Y, to.Y))
	timeToReachPlanet += consts.TimeToReachNearestPlanetS * uint64(cooridanteDiff(from.Z, to.Z))

	return time.Duration(timeToReachPlanet) * time.Second
}

func cooridanteDiff[T coordinate](a, b T) T {
	if a > b {
		return a - b
	}

	return b - a
}
