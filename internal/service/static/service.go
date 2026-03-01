package static

import (
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type registryProvider interface {
	GetAllBuildingStats() ([]registry.BuildingStats, error)
}

type Service struct {
	registry registryProvider
}

func New(registry registryProvider) *Service {
	return &Service{registry: registry}
}
