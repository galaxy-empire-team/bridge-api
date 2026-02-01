package planet

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type systemStorage interface {
	GetSystemPlanets(ctx context.Context, system models.System) (models.SystemPlanets, error)
}

type Service struct {
	systemStorage systemStorage
}

func New(systemStorage systemStorage) *Service {
	return &Service{
		systemStorage: systemStorage,
	}
}
