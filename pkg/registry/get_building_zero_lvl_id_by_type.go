package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBuildingZeroLvlIDByType(buildingType consts.BuildingType) (consts.BuildingID, error) {
	buildingID, exists := r.zeroLvlBuildings[buildingType]
	if !exists {
		return 0, fmt.Errorf("%w: type %s", ErrNotFound, buildingType)
	}

	return buildingID, nil
}
