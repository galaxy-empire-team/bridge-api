package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

const (
	// Public operations operationID equals 0 == is not set.
	// OperationID is used by event manager to guarantee only once planet creation
	colonizeOperationID = 0
	attemptsCount       = 3
)

func (s *Service) ColonizeCapitol(ctx context.Context, userID uuid.UUID) error {
	planetID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("uuid.NewV7(): %w", err)
	}

	for i := 0; i < attemptsCount; i++ {
		position := models.Coordinates{
			X: consts.PlanetPositionX((s.randomGenerator.Uint32() % consts.GalaxyCount) + 1),
			Y: consts.PlanetPositionY((s.randomGenerator.Uint32() % consts.SystemInGalaxyCount) + 1),
		}
		colonizedSystemPlanets, err := s.planetStorage.GetColonizedCoordinates(ctx, position)
		if err != nil {
			return fmt.Errorf("planetStorage.GetColonizedCoordinates(): %w", err)
		}

		availablePositionZ, ok := s.findAvailablePositionZ(colonizedSystemPlanets)
		if !ok {
			s.log.Info("no available position found", zap.Int("attempt", i), zap.Any("position", position))
			continue
		}

		position.Z = availablePositionZ
		planetToColonize := models.Planet{
			ID:          planetID,
			UserID:      userID,
			Coordinates: position,
			IsCapitol:   true,
			HasMoon:     false,
		}

		err = s.planetStorage.ColonizePlanet(ctx, planetToColonize, colonizeOperationID)
		if err != nil {
			return fmt.Errorf("planetStorage.ColonizePlanet(): %w", err)
		}
	}

	return nil
}
