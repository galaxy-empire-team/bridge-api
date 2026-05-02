package planethandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type PlanetService interface {
	CreateCapitol(ctx context.Context, userID uuid.UUID) error
	GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error)
	GetPlanet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Planet, error)
	GetPlanetResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Resources, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	GetFleet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Fleet, error)
	GetResearches(ctx context.Context, userID uuid.UUID) (models.UserResearches, error)
	StartBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) (models.FinishTime, error)
	StartResearch(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, currentResearchID consts.ResearchID) (models.ResearchProgressInfo, error)
	StartFleetConstruction(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, fleet models.FleetUnitCount) (models.FleetUnitConstructionInfo, error)
}
