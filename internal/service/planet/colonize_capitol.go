package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (s *Service) ColonizeCapitol(ctx context.Context, userID uuid.UUID) error {
	planetID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("uuid.NewV7(): %w", err)
	}

	planetToColonize := models.Planet{
		ID: planetID,
		X: uint8(s.randomGenerator.Uint32() % galaxyCount),
		Y: uint8(s.randomGenerator.Uint32() % systemInGalaxyCount),
		Z: uint8(s.randomGenerator.Uint32() % planetsInSystemCount),
	}
	err = s.planetRepo.ColonizeCapitol(ctx, userID, planetToColonize)
	if err != nil {
		return fmt.Errorf("planetRepo.ColonizeCapitol(): %w", err)
	}

	return nil
}
