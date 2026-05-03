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
	ColonizePlanet(ctx context.Context, planet models.Planet, operationID uint64) error
	GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]models.PlanetIDWithCapitol, error)
	GetCoordinates(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error)
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetBuildsInProgressCount(ctx context.Context, planetID uuid.UUID) (uint8, error)
	GetCurrentBuilds(ctx context.Context, planetID uuid.UUID) ([]models.BuildingInProgress, error)
	GetCurrentFleetConstruction(ctx context.Context, planetID uuid.UUID) (models.FleetUnitConstructionInfo, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	GetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.FleetUnitCount, error)
	GetAllPlanetBuildings(ctx context.Context, userID uuid.UUID) ([]consts.BuildingID, error)
	CheckPlanetBelongsToUser(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (bool, error)
	CheckFleetConstruction(ctx context.Context, planetID uuid.UUID) (bool, error)
}

type researchStorage interface {
	GetAllUserResearches(ctx context.Context, userID uuid.UUID) ([]consts.ResearchID, error)
	GetUserResearchesProgress(ctx context.Context, userID uuid.UUID) ([]models.ResearchProgressInfo, error)
	CheckResearchInProgress(ctx context.Context, user_id uuid.UUID) (bool, error)
}

// Separate storage methods that executes inside a transaction
type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	CreateBuildingEvent(ctx context.Context, buildEvent models.BuildEvent) error
	CreateResearchEvent(ctx context.Context, researchEvent models.ResearchEvent) error
	CreateFleetConstructEvent(ctx context.Context, fleetConstructEvent models.FleetConstructEvent) error
}

type txManager interface {
	ExecPlanetTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type resourceRepository interface {
	RecalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
	RecalcResourcesWithUpdatedTime(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedAt time.Time) error
}

type registryProvider interface {
	GetBuildingStatsByID(buildingID consts.BuildingID) (registry.BuildingStats, error)
	GetBuildingNextLvlID(buildingID consts.BuildingID) (consts.BuildingID, error)
	GetResearchNextLvlID(researchID consts.ResearchID) (consts.ResearchID, error)
	GetFleetUnitStatsByID(fleetUnitID consts.FleetUnitID) (registry.FleetUnitStats, error)
	GetResearchStatsByID(researchID consts.ResearchID) (registry.ResearchStats, error)
}

type Service struct {
	txManager          txManager
	planetStorage      planetStorage
	researchStorage    researchStorage
	resourceRepository resourceRepository
	registry           registryProvider
	randomGenerator    *rand.Rand
}

func New(
	txManager txManager,
	planetStorage planetStorage,
	researchStorage researchStorage,
	resourceRepository resourceRepository,
	registry registryProvider,
) *Service {
	gen := rand.New(rand.NewSource((time.Now().UnixNano())))

	return &Service{
		txManager:          txManager,
		planetStorage:      planetStorage,
		researchStorage:    researchStorage,
		resourceRepository: resourceRepository,
		registry:           registry,
		randomGenerator:    gen,
	}
}
