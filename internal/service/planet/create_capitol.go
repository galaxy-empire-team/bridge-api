package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) CreateCapitol(ctx context.Context, userID uuid.UUID) error {
	planetID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("uuid.NewV7(): %w", err)
	}

	generatedLocation := models.Coordinates{
		X: uint8(s.randomGenerator.Uint32() % galaxyCount),
		Y: uint8(s.randomGenerator.Uint32() % systemInGalaxyCount),
		Z: uint8(s.randomGenerator.Uint32() % planetsInSystemCount),
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
		return fmt.Errorf("planetRepo.CreatePlanet(): %w", err)
	}

	return nil
}
