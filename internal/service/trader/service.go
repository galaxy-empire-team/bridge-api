package trader

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type planetStorage interface {
	CheckPlanetHasMoon(ctx context.Context, planetID uuid.UUID) (bool, error)
}

// Separate storage methods that execute inside a transaction.
type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	GetUserResourcesForUpdate(ctx context.Context, userID uuid.UUID) (models.UserResources, error)
	SetUserResources(ctx context.Context, userID uuid.UUID, resources models.UserResources) error
	AddBoost(ctx context.Context, userID uuid.UUID, boost models.UserBoost) error
	AddFleet(ctx context.Context, planetID uuid.UUID, fleetUnit models.FleetUnitCount) error
	AddMistLaunches(ctx context.Context, userID uuid.UUID, count uint64) error
	SetHasMoon(ctx context.Context, planetID uuid.UUID, hasMoon bool) error
	CheckPlanetHasMoon(ctx context.Context, planetID uuid.UUID) (bool, error)
}

type txManager interface {
	ExecTraderTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type repository interface {
	CheckPlanetOwner(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
	RecalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
}

type registryProvider interface {
	GetTraderItemStatsByID(id consts.TraderItemID) (registry.TraderItemStats, error)
}

type Service struct {
	txManager     txManager
	planetStorage planetStorage
	repository    repository
	registry      registryProvider
}

func New(
	txManager txManager,
	planetStorage planetStorage,
	repository repository,
	registry registryProvider,
) *Service {
	return &Service{
		txManager:     txManager,
		planetStorage: planetStorage,
		repository:    repository,
		registry:      registry,
	}
}
