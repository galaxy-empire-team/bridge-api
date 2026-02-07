package planet

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type planetStorage interface {
	CreatePlanet(ctx context.Context, planet models.Planet) error
	GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]models.PlanetIDWithCapitol, error)
	GetResources(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetCoordinates(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error)
	GetBuildingsInfo(ctx context.Context, planetID uuid.UUID, BuildingTypes []consts.BuildingType) (map[consts.BuildingType]models.BuildingInfo, error)
	GetCurrentBuildsCount(ctx context.Context, planetID uuid.UUID) (uint8, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
}

// Separate storage methods that executes inside a transaction
type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	GetBuildingID(ctx context.Context, planetID uuid.UUID, BuildingType consts.BuildingType) (consts.BuildingID, error)
	GetBuildingsInfo(ctx context.Context, planetID uuid.UUID, BuildingTypes []consts.BuildingType) (map[consts.BuildingType]models.BuildingInfo, error)
	CreateBuildingEvent(ctx context.Context, buildEvent models.BuildEvent) error
	SetFinishedBuildingTime(ctx context.Context, planetID uuid.UUID, buildinID consts.BuildingID, finishedAt time.Time) error
	CreateBuilding(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID) error
}

type txManager interface {
	ExecPlanetTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type registryProvider interface {
	GetBuildingZeroLvlStats(buildingType consts.BuildingType) (registry.BuildingStats, error)
	GetBuildingStatsByID(buildingID consts.BuildingID) (registry.BuildingStats, error)
	GetBuildingStatsByType(buildingType consts.BuildingType, level consts.BuildingLevel) (registry.BuildingStats, error)
	GetBuildingNextLvlStats(buildingID consts.BuildingID) (registry.BuildingStats, error)
}

type Service struct {
	txManager       txManager
	planetStorage   planetStorage
	registry        registryProvider
	randomGenerator *rand.Rand
}

func New(txManager txManager, planetStorage planetStorage, registry registryProvider) *Service {
	gen := rand.New(rand.NewSource((time.Now().UnixNano())))

	return &Service{
		txManager:       txManager,
		planetStorage:   planetStorage,
		registry:        registry,
		randomGenerator: gen,
	}
}
