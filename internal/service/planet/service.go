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
	GetCoordinates(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error)
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetPlanetMinesProduction(ctx context.Context, planetID uuid.UUID) (map[consts.BuildingType]uint64, error)
	GetBuildsInProgressCount(ctx context.Context, planetID uuid.UUID) (uint8, error)
	GetCurrentBuilds(ctx context.Context, planetID uuid.UUID) ([]models.BuildingInProgress, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	GetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.PlanetFleetUnitCount, error)
	CheckPlanetBelongsToUser(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (bool, error)
	GetAllPlanetBuildings(ctx context.Context, userID uuid.UUID) ([]consts.BuildingID, error)
}

type researchStorage interface {
	GetUserResearches(ctx context.Context, userID uuid.UUID) ([]consts.ResearchID, error)
}

// Separate storage methods that executes inside a transaction
type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	CreateBuildingEvent(ctx context.Context, buildEvent models.BuildEvent) error
	SetFinishedBuildingTime(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID, finishedAt time.Time) error
}

type txManager interface {
	ExecPlanetTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type registryProvider interface {
	GetBuildingStatsByID(buildingID consts.BuildingID) (registry.BuildingStats, error)
	GetBuildingNextLvlID(buildingID consts.BuildingID) (consts.BuildingID, error)
	GetFleetUnitStatsByID(fleetUnitID consts.FleetUnitID) (registry.FleetUnitStats, error)
	GetResearchStatsByID(researchID consts.ResearchID) (registry.ResearchStats, error)
}

type Service struct {
	txManager       txManager
	planetStorage   planetStorage
	researchStorage researchStorage
	registry        registryProvider
	randomGenerator *rand.Rand
}

func New(txManager txManager, planetStorage planetStorage, researchStorage researchStorage, registry registryProvider) *Service {
	gen := rand.New(rand.NewSource((time.Now().UnixNano())))

	return &Service{
		txManager:       txManager,
		planetStorage:   planetStorage,
		researchStorage: researchStorage,
		registry:        registry,
		randomGenerator: gen,
	}
}
