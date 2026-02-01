package planet

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type planetStorage interface {
	CheckPlanetExists(ctx context.Context, coordinates models.Coordinates) (bool, error)
	CheckPlanetBelongsToUser(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (bool, error)
}

type TxStorages interface {
	// --- missionStorage ---
	CreateMissionEvent(ctx context.Context, colonizeEvent models.MissionEvent) error
}

type txManager interface {
	ExecMissionTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type Service struct {
	txManager     txManager
	planetStorage planetStorage
}

func New(txManager txManager, planetStorage planetStorage) *Service {
	return &Service{
		txManager:     txManager,
		planetStorage: planetStorage,
	}
}
