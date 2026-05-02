package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

// ColonizePlanet allows a user to colonize a new planet.
// OperationID is optional and used by the event manager to make a guarantee of at-least-once execution.
func (s *Service) ColonizePlanet(ctx context.Context, userID uuid.UUID, req CreatePlanetRequest) error {
	planetID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("uuid.NewV7(): %w", err)
	}

	planetToColonize := models.Planet{
		ID:          planetID,
		UserID:      userID,
		Coordinates: req.Coordinates,
		Resources:   req.Resources,
		IsCapitol:   req.IsCapitol,
		HasMoon:     false,
	}

	err = s.planetStorage.ColonizePlanet(ctx, planetToColonize, req.OperationID)
	if err != nil {
		return fmt.Errorf("planetStorage.ColonizePlanet(): %w", err)
	}

	return nil
}
