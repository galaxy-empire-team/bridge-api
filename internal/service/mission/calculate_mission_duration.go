package mission

import (
	"fmt"
	"time"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type coordinate interface {
	consts.PlanetPositionX | consts.PlanetPositionY | consts.PlanetPositionZ
}

func (s *Service) calculateMissionDuration(planetFrom models.Coordinates, planetTo models.Coordinates, fleet []models.FleetUnitCount, speedMultiplier float64) (time.Duration, error) {
	if speedMultiplier < consts.SpeedMultiplierMin || speedMultiplier > consts.SpeedMultiplierMax {
		return 0, fmt.Errorf("invalid speed multiplier: %f", speedMultiplier)
	}

	minSpeed, err := s.calcMinSpeed(fleet)
	if err != nil {
		return 0, fmt.Errorf("calcMinSpeed(): %w", err)
	}

	duration := calcDistanceDuration(planetFrom, planetTo) / time.Duration(minSpeed)
	duration = time.Duration(float64(duration) / speedMultiplier)

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
