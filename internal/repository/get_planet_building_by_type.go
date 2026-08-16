package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

// GetBuildingByType is a wrapper on GetBuildingsByTypes to handle only one building type.
// Method returns stats for the target research type.
// If tech is not found in the database it returns zero-lvl stats.
func (r *Repository) GetBuildingByType(ctx context.Context, planetID uuid.UUID, buildingType consts.BuildingType) (registry.BuildingStats, error) {
	buildings, err := r.GetBuildingsByTypes(ctx, planetID, []consts.BuildingType{buildingType})
	if err != nil {
		return registry.BuildingStats{}, fmt.Errorf("r.GetBuildingsByTypes(): %w", err)
	}

	buildingStats, ok := buildings[buildingType]
	if !ok {
		return registry.BuildingStats{}, fmt.Errorf("%w: %s", ErrInvalidBuildingType, buildingType)
	}

	return buildingStats, nil
}

// GetBuildingsByTypes returns stats for building types.
// If building is not found in the database it returns zero-lvl stats.
func (r *Repository) GetBuildingsByTypes(ctx context.Context, planetID uuid.UUID, buildingTypes []consts.BuildingType) (map[consts.BuildingType]registry.BuildingStats, error) {
	buildingTypeToID, err := r.planetStorage.GetPlanetBuildingsByTypes(ctx, planetID, buildingTypes)
	if err != nil {
		return nil, fmt.Errorf("planetStorage.GetPlanetBuildingsByTypes(): %w", err)
	}

	res := make(map[consts.BuildingType]registry.BuildingStats)
	for _, buildingType := range buildingTypes {
		buildingID, ok := buildingTypeToID[buildingType]
		if !ok {
			buildingID, err = r.registry.GetBuildingZeroLvlIDByType(buildingType)
			if err != nil {
				return nil, fmt.Errorf("registry.GetBuildingZeroLvlIDByType(%s): %w", buildingType, err)
			}
		}

		buildingStats, err := r.registry.GetBuildingStatsByID(buildingID)
		if err != nil {
			return nil, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
		}

		res[buildingType] = buildingStats
	}

	return res, nil
}
