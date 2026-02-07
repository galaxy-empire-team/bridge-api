package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBuildingStatsByID(buildingID consts.BuildingID) (BuildingStats, error) {
	stat, ok := r.buildings[buildingID]
	if !ok {
		return BuildingStats{}, fmt.Errorf("%w: ID %d", ErrNotFound, buildingID)
	}

	return stat, nil
}
