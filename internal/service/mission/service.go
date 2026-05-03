package mission

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

type planetStorage interface {
	GetIDByCoordinates(ctx context.Context, coordinates models.Coordinates) (uuid.UUID, error)
	CheckPlanetExists(ctx context.Context, coordinates models.Coordinates) (bool, error)
	CheckPlanetBelongsToUser(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (bool, error)
	GetUserPlanetsCount(ctx context.Context, userID uuid.UUID) (uint8, error)
	GetCoordinates(ctx context.Context, planetID uuid.UUID) (models.Coordinates, error)
}

type missionStorage interface {
	GetCurrentUserMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error)
}

type researchStorage interface {
	GetUserResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType) (map[consts.ResearchType]consts.ResearchID, error)
}

type TxStorages interface {
	// --- planetStorage ---
	GetFleetForUpdate(ctx context.Context, planetID uuid.UUID) ([]models.FleetUnitCount, error)
	SetFleet(ctx context.Context, planetID uuid.UUID, fleet []models.FleetUnitCount) error
	GetResourcesForUpdate(ctx context.Context, planetID uuid.UUID) (models.Resources, error)
	SetResources(ctx context.Context, planetID uuid.UUID, updatedResources models.Resources) error
	// --- missionStorage ---
	CreateMissionEvent(ctx context.Context, colonizeEvent models.MissionEvent) error
}

type txManager interface {
	ExecMissionTx(ctx context.Context, fn func(ctx context.Context, storages TxStorages) error) error
}

type resourceRepository interface {
	RecalcResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
	RecalcResourcesWithUpdatedTime(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedAt time.Time) error
}

type registryProvider interface {
	CheckFleetUnitIDExists(fleetUnitID consts.FleetUnitID) bool
	GetFleetUnitTypeCount() int
	GetMissionIDByType(missionType consts.MissionType) (consts.MissionID, error)
	GetFleetUnitStatsByID(fleetUnitID consts.FleetUnitID) (registry.FleetUnitStats, error)
	GetResearchStatsByID(researchID consts.ResearchID) (registry.ResearchStats, error)
	GetFleetUnitIDs() []consts.FleetUnitID
}

type Service struct {
	txManager          txManager
	researchStorage    researchStorage
	planetStorage      planetStorage
	missionStorage     missionStorage
	resourceRepository resourceRepository
	registry           registryProvider
}

func New(
	txManager txManager,
	planetStorage planetStorage,
	missionStorage missionStorage,
	researchStorage researchStorage,
	resourceRepository resourceRepository,
	registry registryProvider,
) *Service {
	return &Service{
		txManager:          txManager,
		researchStorage:    researchStorage,
		planetStorage:      planetStorage,
		missionStorage:     missionStorage,
		resourceRepository: resourceRepository,
		registry:           registry,
	}
}
