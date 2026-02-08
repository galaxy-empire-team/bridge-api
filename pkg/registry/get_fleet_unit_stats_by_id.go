package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetFleetUnitStatsByID(fleetUnitID consts.FleetUnitID) (FleetUnitStats, error) {
	stat, ok := r.fleet[fleetUnitID]
	if !ok {
		return FleetUnitStats{}, fmt.Errorf("%w: ID %d", ErrNotFound, fleetUnitID)
	}

	return stat, nil
}
