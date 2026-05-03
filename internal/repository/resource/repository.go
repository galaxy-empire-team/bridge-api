package planet

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type researchStorage interface {
	GetUserResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType) (map[consts.ResearchType]consts.ResearchID, error)
}

type planetStorage interface {
	GetPlanetMinesProduction(ctx context.Context, planetID uuid.UUID) (map[consts.BuildingType]uint64, error)
}

// Separate storage methods that executes inside a transaction
type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
}

type txManager interface {
	ExecResourceRepoTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type registryProvider interface {
	GetResearchStatsByID(researchID consts.ResearchID) (registry.ResearchStats, error)
}

type Repository struct {
	txManager       txManager
	planetStorage   planetStorage
	researchStorage researchStorage
	registry        registryProvider
}

func New(txManager txManager, planetStorage planetStorage, researchStorage researchStorage, registry registryProvider) *Repository {
	return &Repository{
		txManager:       txManager,
		planetStorage:   planetStorage,
		researchStorage: researchStorage,
		registry:        registry,
	}
}
