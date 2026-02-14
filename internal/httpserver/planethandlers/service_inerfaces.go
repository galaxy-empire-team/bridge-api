package planethandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

type PlanetService interface {
	GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error)
	CreateCapitol(ctx context.Context, userID uuid.UUID) error
	UpgradeBuilding(ctx context.Context, planetID uuid.UUID, buildingType string) error
	GetBuildingStats(ctx context.Context, buildingType string, level uint8) (models.BuildingStats, error)
	GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error)
	GetPlanet(ctx context.Context, planetID uuid.UUID) (models.Planet, error)
	GetFleet(ctx context.Context, planetID uuid.UUID) ([]models.PlanetFleetUnitCount, error)
}
