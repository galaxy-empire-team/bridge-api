package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) updatePlanetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.PlanetFleetUnitCount, storage TxStorages) error {
	reqFleet := lo.Associate(fleet, func(fleetUnit models.PlanetFleetUnitCount) (consts.FleetUnitID, uint64) {
		return fleetUnit.ID, fleetUnit.Count
	})

	planetFleet, err := storage.GetFleetForUpdate(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetFleetCountForUpdate(): %w", err)
	}

	planetFleetMap := lo.Associate(planetFleet, func(fleetUnit models.PlanetFleetUnitCount) (consts.FleetUnitID, uint64) {
		return fleetUnit.ID, fleetUnit.Count
	})

	var leftFleetUnits []models.PlanetFleetUnitCount
	for _, fleetUnitID := range s.registry.GetFleetUnitIDs() {
		reqCount, ok := reqFleet[fleetUnitID]
		if !ok {
			continue
		}

		planetCount, ok := planetFleetMap[fleetUnitID]
		if !ok {
			return models.ErrFleetNotFound
		}

		if planetCount < reqCount {
			return models.ErrNotEnoughFleetUnits
		}

		leftFleetUnits = append(leftFleetUnits, models.PlanetFleetUnitCount{
			ID:    fleetUnitID,
			Count: planetCount - reqCount,
		})
	}

	err = storage.SetFleet(ctx, planetID, leftFleetUnits)
	if err != nil {
		return fmt.Errorf("planetStorage.SetFleet(): %w", err)
	}

	return nil
}
