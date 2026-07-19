package planet

import (
	"slices"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) findAvailablePositionZ(colonizedSystemPlanets []consts.PlanetPositionZ) (consts.PlanetPositionZ, bool) {
	initialPosZ := consts.PlanetPositionZ(s.randomGenerator.Uint32()%consts.PlanetsInSystemCount + 1)

	shouldSkipPosZ := func(posZ consts.PlanetPositionZ) bool {
		if posZ == consts.NPCTierOnePositionZ || posZ == consts.NPCTierTwoPositionZ || posZ == consts.NPCTierThreePositionZ {
			return true
		}

		if slices.Contains(colonizedSystemPlanets, posZ) {
			return true
		}

		return false
	}

	for i := initialPosZ; i <= consts.PlanetsInSystemCount; i++ {
		if shouldSkipPosZ(i) {
			continue
		}

		return i, true
	}

	for i := consts.PlanetPositionZ(1); i < initialPosZ; i++ {
		if shouldSkipPosZ(i) {
			continue
		}

		return i, true
	}

	return 0, false
}
