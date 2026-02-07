package planet

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) GetBuildingStats(ctx context.Context, buildingType string, level uint8) (models.BuildingStats, error) {
	if !consts.IsValidBuildingType(consts.BuildingType(buildingType)) {
		return models.BuildingStats{}, models.ErrBuildTypeInvalid
	}

	if level > consts.MaxBuildingLvl {
		return models.BuildingStats{}, models.ErrBuildingInvalidLevel
	}

	buildingStats, err := s.registry.GetBuildingStatsByType(consts.BuildingType(buildingType), consts.BuildingLevel(level))
	if err != nil {
		return models.BuildingStats{}, fmt.Errorf("planetRepo.GetBuildingStats(): %w", err)
	}

	return toModelBuildingStats(buildingStats), nil
}
