package planethandlers

import (
	"context"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

type PlanetService interface {
	GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error)
	CreateCapitol(ctx context.Context, userID uuid.UUID) error
	UpgradeBuilding(ctx context.Context, planetID uuid.UUID, buildingType string) error
}
