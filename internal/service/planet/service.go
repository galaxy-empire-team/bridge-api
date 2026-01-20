package planet

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

const (
	galaxyCount          = 1
	systemInGalaxyCount  = 3
	planetsInSystemCount = 16

	defaultLvl = 0
)

type planetStorage interface {
	GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error)
	CreateCapitol(ctx context.Context, userID uuid.UUID, planet models.Planet) error
	GetUserPlanetIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SaveResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	GetBuildStats(ctx context.Context, buildType models.BuildType, level uint8) (models.BuildStats, error)
	GetBuildLvl(ctx context.Context, planetID uuid.UUID, buildType models.BuildType) (uint8, error)
	GetBuildsInfo(ctx context.Context, planetID uuid.UUID, buildTypes []models.BuildType) (map[models.BuildType]models.PlanetBuildInfo, error)
	SaveBuild(ctx context.Context, planetID uuid.UUID, buildInfo models.PlanetBuildInfo) error
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
