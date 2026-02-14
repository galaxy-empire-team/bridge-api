package registry

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) CheckFleetUnitIDExists(fleetUnitID consts.FleetUnitID) bool {
	_, ok := r.fleet[fleetUnitID]
	return ok
}
