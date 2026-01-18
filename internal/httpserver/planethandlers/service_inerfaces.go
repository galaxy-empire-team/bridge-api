package planethandlers

import (
	"context"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

type PlanetService interface {
	GetCapitolPlanet(ctx context.Context, userID uuid.UUID) (models.Planet, error)
	ColonizeCapitol(ctx context.Context, userID uuid.UUID) error
}
