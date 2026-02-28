package registry

import "github.com/galaxy-empire-team/bridge-api/pkg/consts"

func (r *Registry) GetFleetUnitIDs() []consts.FleetUnitID {
	ids := make([]consts.FleetUnitID, 0, len(r.fleet))
	for id := range r.fleet {
		ids = append(ids, id)
	}

	return ids
}
