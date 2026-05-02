package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

const (
	// Public operations operationID equals 0 == is not set.
	// OperationID is used by event manager to guarantee only once planet creation
	colonizeOperationID = 0
)

func (s *Service) ColonizeCapitol(ctx context.Context, userID uuid.UUID) error {
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

	err = s.planetStorage.ColonizePlanet(ctx, planetToColonize, colonizeOperationID)
	if err != nil {
		return fmt.Errorf("planetStorage.ColonizePlanet(): %w", err)
	}

	return nil
}
