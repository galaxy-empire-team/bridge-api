package planet

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

const (
	galaxyCount          = 1
	systemInGalaxyCount  = 3
	planetsInSystemCount = 16

	defaultLvl         = 0
	defaultProductionS = 10
	maxBuildingLvl     = 5

	maxBuildingsInProgress = 2
)

type planetStorage interface {
	CreatePlanet(ctx context.Context, planet models.Planet) error
	GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]models.PlanetIDWithCapitol, error)
	GetResources(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetCoordinates(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error)
	GetBuildingsInfo(ctx context.Context, planetID uuid.UUID, BuildingTypes []consts.BuildingType) (map[consts.BuildingType]models.BuildingInfo, error)
	GetBuildingStats(ctx context.Context, BuildingType consts.BuildingType, level uint8) (models.BuildingStats, error)
	GetCurrentBuildsCount(ctx context.Context, planetID uuid.UUID) (uint8, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
}

// Separate storage methods that executes inside a transaction
type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	GetBuildingStats(ctx context.Context, BuildingType consts.BuildingType, level uint8) (models.BuildingStats, error)
	GetBuildingLvl(ctx context.Context, planetID uuid.UUID, BuildingType consts.BuildingType) (uint8, error)
	GetBuildingsInfo(ctx context.Context, planetID uuid.UUID, BuildingTypes []consts.BuildingType) (map[consts.BuildingType]models.BuildingInfo, error)
	CreateBuildingEvent(ctx context.Context, buildEvent models.BuildEvent) error
	SetFinishedBuildingTime(ctx context.Context, planetID uuid.UUID, buildingInfo models.BuildingInfo) error
}

type txManager interface {
	ExecPlanetTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type Service struct {
	txManager       txManager
	planetStorage   planetStorage
	randomGenerator *rand.Rand
}

func New(txManager txManager, planetStorage planetStorage) *Service {
	gen := rand.New(rand.NewSource((time.Now().UnixNano())))

	return &Service{
		txManager:       txManager,
		planetStorage:   planetStorage,
		randomGenerator: gen,
	}
}
