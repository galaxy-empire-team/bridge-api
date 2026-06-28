package event

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type eventStorage interface {
	GetBuildsInProgressCount(ctx context.Context, planetID uuid.UUID) (uint8, error)
	CheckFleetConstruction(ctx context.Context, planetID uuid.UUID) (bool, error)
	CheckResearchInProgress(ctx context.Context, userID uuid.UUID) (bool, error)
}

type planetStorage interface {
	GetAllPlanetBuildings(ctx context.Context, planetID uuid.UUID) ([]consts.BuildingID, error)
}

type researchStorage interface {
	GetAllUserResearches(ctx context.Context, userID uuid.UUID) ([]consts.ResearchID, error)
}

// Separate storage methods that executes inside a transaction
type TxStorages interface {
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	AddResources(ctx context.Context, planetID uuid.UUID, resources models.Resources) error
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	CreateBuildingEvent(ctx context.Context, buildEvent models.BuildEvent) error
	CreateResearchEvent(ctx context.Context, researchEvent models.ResearchEvent) error
	CreateFleetConstructEvent(ctx context.Context, fleetConstructEvent models.FleetConstructEvent) error
	DeleteBuildingEvent(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID) error
	DeleteResearchEvent(ctx context.Context, userID uuid.UUID, researchID consts.ResearchID) error
	DeleteFleetConstructionEvent(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	GetBoostByIDForUpdate(ctx context.Context, userID uuid.UUID, boostID consts.BoostID) (models.UserBoost, error)
	SetBoost(ctx context.Context, userID uuid.UUID, boost models.UserBoost) error
	GetResearchEventForUpdate(ctx context.Context, userID uuid.UUID, researchID consts.ResearchID) (models.EventFinishTime, error)
	GetBuildingEventForUpdate(ctx context.Context, planetID uuid.UUID, buildingID consts.BuildingID) (models.EventFinishTime, error)
	GetFleetConstructionEventForUpdate(ctx context.Context, planetID uuid.UUID) (models.EventFinishTime, error)
	SetResearchFinishTime(ctx context.Context, researchEvent models.EventFinishTime) error
	SetBuildingFinishTime(ctx context.Context, buildingEvent models.EventFinishTime) error
	SetFleetConstructionFinishTime(ctx context.Context, fleetConstructionEvent models.EventFinishTime) error
}

type txManager interface {
	ExecEventTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type repository interface {
	CheckPlanetOwner(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
	RecalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
}

type registryProvider interface {
	GetBuildingStatsByID(buildingID consts.BuildingID) (registry.BuildingStats, error)
	GetBuildingNextLvlID(buildingID consts.BuildingID) (consts.BuildingID, error)
	GetResearchNextLvlID(researchID consts.ResearchID) (consts.ResearchID, error)
	GetFleetUnitStatsByID(fleetUnitID consts.FleetUnitID) (registry.FleetUnitStats, error)
	GetResearchStatsByID(researchID consts.ResearchID) (registry.ResearchStats, error)
	GetBoostStatsByID(boostID consts.BoostID) (registry.BoostStats, error)
}

type Service struct {
	txManager       txManager
	eventStorage    eventStorage
	planetStorage   planetStorage
	researchStorage researchStorage
	repository      repository
	registry        registryProvider
}

func New(
	txManager txManager,
	eventStorage eventStorage,
	planetStorage planetStorage,
	researchStorage researchStorage,
	repository repository,
	registry registryProvider,
) *Service {
	return &Service{
		txManager:       txManager,
		eventStorage:    eventStorage,
		planetStorage:   planetStorage,
		researchStorage: researchStorage,
		repository:      repository,
		registry:        registry,
	}
}
