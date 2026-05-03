package mission

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) checkFleetValid(fleet []models.FleetUnitCount) error {
	if len(fleet) == 0 {
		return models.ErrFleetCannotBeEmpty
	}

	if len(fleet) > s.registry.GetFleetUnitTypeCount() {
		return models.ErrInvalidInput
	}

	duplicateFleetUnitIDs := make(map[consts.FleetUnitID]struct{})
	for _, fleetUnit := range fleet {
		if fleetUnit.Count == 0 {
			return models.ErrFleetUnitCountCannotBeZero
		}

		if !s.registry.CheckFleetUnitIDExists(fleetUnit.ID) {
			return fmt.Errorf("%w: %d", models.ErrFleetIDNotExists, fleetUnit.ID)
		}

		if _, exists := duplicateFleetUnitIDs[fleetUnit.ID]; exists {
			return fmt.Errorf("%w: %d", models.ErrDuplicateFleetUnitID, fleetUnit.ID)
		}

		duplicateFleetUnitIDs[fleetUnit.ID] = struct{}{}
	}

	return nil
}
