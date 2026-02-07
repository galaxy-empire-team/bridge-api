package registry

import (
	"fmt"
	"slices"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBuildingStatsByType(buildingType consts.BuildingType, level consts.BuildingLevel) (BuildingStats, error) {
	if level > consts.MaxBuildingLvl {
		return BuildingStats{}, fmt.Errorf("%w: level %d", ErrInvalidBuildingLevel, level)
	}

	if !slices.Contains(consts.GetBuildingTypes(), buildingType) {
		return BuildingStats{}, fmt.Errorf("%w: type %s", ErrInvalidBuildingType, buildingType)
	}

	for _, stat := range r.buildings {
		if stat.Type == buildingType && stat.Level == level {
			return stat, nil
		}
	}

	return BuildingStats{}, fmt.Errorf("%w: type %s level %d", ErrNotFound, buildingType, level)
}
