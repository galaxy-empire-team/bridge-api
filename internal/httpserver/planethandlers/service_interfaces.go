package planethandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type PlanetService interface {
	GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error)
	GetPlanet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Planet, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	CreateCapitol(ctx context.Context, userID uuid.UUID) error
	UpgradeBuilding(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, buildingID consts.BuildingID) error
	GetFleet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) ([]models.PlanetFleetUnitCount, error)
}
