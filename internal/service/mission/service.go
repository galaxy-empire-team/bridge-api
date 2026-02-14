package planet

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type planetStorage interface {
	CheckPlanetExists(ctx context.Context, coordinates models.Coordinates) (bool, error)
	CheckPlanetBelongsToUser(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (bool, error)
}

type TxStorages interface {
	// --- planetStorage ---
	GetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.PlanetFleetUnitCount, error)
	SetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.PlanetFleetUnitCount) error
	// --- missionStorage ---
	CreateMissionEvent(ctx context.Context, colonizeEvent models.MissionEvent) error
}

type txManager interface {
	ExecMissionTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type registryProvider interface {
	CheckFleetUnitIDExists(fleetUnitID consts.FleetUnitID) bool
	GetFleetUnitTypeCount() int
}

type Service struct {
	txManager     txManager
	planetStorage planetStorage
	registry      registryProvider
}

func New(txManager txManager, planetStorage planetStorage, registry registryProvider) *Service {
	return &Service{
		txManager:     txManager,
		planetStorage: planetStorage,
		registry:      registry,
	}
}
