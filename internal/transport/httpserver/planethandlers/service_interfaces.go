package planethandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	planetservice "github.com/galaxy-empire-team/bridge-api/internal/service/planet"
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
	GetUserResources(ctx context.Context, userID uuid.UUID) (planetservice.GetUserResourcesResponse, error)
}
