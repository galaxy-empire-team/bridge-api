package planethandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type PlanetService interface {
	ColonizeCapitol(ctx context.Context, userID uuid.UUID) error
	GetCapitolID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetPlanet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Planet, error)
	GetPlanetResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Resources, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	GetFleet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Fleet, error)
	GetResearches(ctx context.Context, userID uuid.UUID) (models.UserResearches, error)
	GetBuildings(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Buildings, error)
	StartBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) (models.FinishTime, error)
	StartResearch(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, currentResearchID consts.ResearchID) (models.ResearchProgressInfo, error)
	StartFleetConstruction(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, fleet models.FleetUnitCount) (models.FleetUnitConstructionInfo, error)
	CancelBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) error
	CancelResearch(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, researchID consts.ResearchID) error
	CancelFleetConstruction(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) error
}
