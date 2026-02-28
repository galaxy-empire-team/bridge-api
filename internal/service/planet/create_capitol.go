package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) CreateCapitol(ctx context.Context, userID uuid.UUID) error {
	planetID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("uuid.NewV7(): %w", err)
	}

	generatedLocation := models.Coordinates{
		X: consts.PlanetPositionX(s.randomGenerator.Uint32() % consts.GalaxyCount),
		Y: consts.PlanetPositionY(s.randomGenerator.Uint32() % consts.SystemInGalaxyCount),
		Z: consts.PlanetPositionZ(s.randomGenerator.Uint32() % consts.PlanetsInSystemCount),
	}
	planetToColonize := models.Planet{
		ID:          planetID,
		UserID:      userID,
		Coordinates: generatedLocation,
		IsCapitol:   true,
		HasMoon:     false,
	}

	err = s.planetStorage.CreatePlanet(ctx, planetToColonize)
	if err != nil {
		return fmt.Errorf("planetStorage.CreatePlanet(): %w", err)
	}

	return nil
}
