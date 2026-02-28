package registry

import (
	"fmt"
	"slices"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBuildingZeroLvlStats(buildingType consts.BuildingType) (BuildingStats, error) {
	if !slices.Contains(consts.GetBuildingTypes(), buildingType) {
		return BuildingStats{}, fmt.Errorf("%w: type %s", ErrInvalidBuildingType, buildingType)
	}

	for _, stat := range r.buildings {
		if stat.Type == buildingType && stat.Level == 1 {
			return stat, nil
		}
	}

	return BuildingStats{}, fmt.Errorf("%w: type %s", ErrNotFound, buildingType)
}
