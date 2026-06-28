package planet

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type planetStorage interface {
	GetCapitolID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetPlanet(ctx context.Context, planetID uuid.UUID) (models.Planet, error)
	GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]models.PlanetIDCapitol, error)
	GetCoordinates(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error)
	GetColonizedCoordinates(ctx context.Context, system models.Coordinates) ([]consts.PlanetPositionZ, error)
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetCurrentBuilds(ctx context.Context, planetID uuid.UUID) ([]models.BuildingInProgress, error)
	GetCurrentFleetConstruction(ctx context.Context, planetID uuid.UUID) (models.FleetUnitConstructionInfo, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	GetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.FleetUnitCount, error)
	GetAllPlanetBuildings(ctx context.Context, planetID uuid.UUID) ([]consts.BuildingID, error)
	CheckPlanetBelongsToUser(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (bool, error)
	ColonizePlanet(ctx context.Context, planet models.Planet, operationID uint64) error
	GetUserResources(ctx context.Context, userID uuid.UUID) (models.UserResources, error)
	GetUserBoosts(ctx context.Context, userID uuid.UUID) ([]models.UserBoost, error)
}

type researchStorage interface {
	GetAllUserResearches(ctx context.Context, userID uuid.UUID) ([]consts.ResearchID, error)
	GetUserResearchesProgress(ctx context.Context, userID uuid.UUID) ([]models.ResearchProgressInfo, error)
}

type repository interface {
	CheckPlanetOwner(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
	RecalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
	RecalcResourcesWithUpdatedTime(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedAt time.Time) error
}

//go:generate mockery --name=randGenerator --filename=rand_generator.go --exported --with-expecter
type randGenerator interface {
	Uint32() uint32
}

type Service struct {
	planetStorage   planetStorage
	researchStorage researchStorage
	repository      repository
	randomGenerator randGenerator
	log             *zap.Logger
}

func New(
	planetStorage planetStorage,
	researchStorage researchStorage,
	repository repository,
	log *zap.Logger,
) *Service {
	gen := rand.New(rand.NewSource((time.Now().UnixNano())))

	return &Service{
		planetStorage:   planetStorage,
		researchStorage: researchStorage,
		repository:      repository,
		randomGenerator: gen,
		log:             log,
	}
}
