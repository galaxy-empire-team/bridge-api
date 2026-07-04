package planet

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type systemStorage interface {
	GetSystemPlanets(ctx context.Context, system models.System) ([]models.PlanetInfo, error)
}

type planetStorage interface {
	GetUserNPCAttacks(ctx context.Context, userID uuid.UUID) ([]models.NPCAttack, error)
}

type Service struct {
	systemStorage systemStorage
	planetStorage planetStorage
}

func New(planetStorage planetStorage, systemStorage systemStorage) *Service {
	return &Service{
		planetStorage: planetStorage,
		systemStorage: systemStorage,
	}
}
