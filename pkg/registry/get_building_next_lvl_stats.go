package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBuildingNextLvlStats(buildingID consts.BuildingID) (BuildingStats, error) {
	_, ok := r.buildings[buildingID]
	if !ok {
		return BuildingStats{}, fmt.Errorf("%w: ID %d", ErrNotFound, buildingID)
	}

	// For simplicity, I assume that the next level of a building has an ID that is exactly 1 greater than the current building's ID.
	// TODO rewrite using map and linked list
	updatedStat, ok := r.buildings[buildingID+1]
	if !ok {
		return BuildingStats{}, fmt.Errorf("%w: ID %d", ErrNotFound, buildingID+1)
	}

	stats, ok := r.buildings[updatedStat.ID]
	if !ok {
		return BuildingStats{}, fmt.Errorf("%w: ID %d", ErrNotFound, updatedStat.ID)
	}

	return stats, nil
}
