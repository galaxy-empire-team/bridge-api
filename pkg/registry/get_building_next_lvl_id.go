package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBuildingNextLvlID(buildingID consts.BuildingID) (consts.BuildingID, error) {
	currentBuilding, ok := r.buildings[buildingID]
	if !ok {
		return 0, fmt.Errorf("%w: ID %d", ErrNotFound, buildingID)
	}

	if currentBuilding.Level >= consts.MaxBuildingLvl {
		return 0, ErrMaxLevelReached
	}

	// For simplicity, I assume that the next level of a building has an ID that is exactly 1 greater than the current building's ID.
	// TODO rewrite using map and linked list
	_, ok = r.buildings[buildingID+1]
	if !ok {
		return 0, fmt.Errorf("%w: ID %d", ErrNotFound, buildingID+1)
	}

	return buildingID + 1, nil
}
