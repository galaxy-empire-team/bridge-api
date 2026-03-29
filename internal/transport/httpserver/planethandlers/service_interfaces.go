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
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	GetFleet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) ([]models.FleetUnitCount, error)
	StartBuildingUpgrade(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) error
	StartResearch(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, currentResearchID consts.ResearchID) error
	StartFleetConstruction(ctx context.Context, userID uuid.UUID, currentPlanet uuid.UUID, fleet models.FleetUnitCount) error
}
