package planet

import (
	"context"
	"fmt"

	"initialservice/internal/models"
)

func (s *Service) GetBuildingStats(ctx context.Context, buildingType string, level uint8) (models.BuildingStats, error) {
	buildingStats, err := s.planetStorage.GetBuildingStats(ctx, models.BuildingType(buildingType), level)
	if err != nil {
		return models.BuildingStats{}, fmt.Errorf("planetRepo.GetBuildingStats(): %w", err)
	}

	return buildingStats, nil
}
