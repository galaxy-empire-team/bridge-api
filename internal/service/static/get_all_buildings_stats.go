package static

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func (s *Service) GetBuildingStats(ctx context.Context) ([]registry.BuildingStats, error) {
	return s.registry.GetAllBuildingStats()
}
