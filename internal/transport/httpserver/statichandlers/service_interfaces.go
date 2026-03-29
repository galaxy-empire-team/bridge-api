package statichandlers

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type StaticService interface {
	GetBuildingStats(ctx context.Context) ([]registry.BuildingStats, error)
}
