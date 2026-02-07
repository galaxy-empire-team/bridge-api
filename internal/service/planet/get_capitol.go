package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error) {
	planetIDs, err := s.planetStorage.GetUserPlanetIDs(ctx, userID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetUserPlanetIDs(): %w", err)
	}

	var capitolID uuid.UUID
	for _, pid := range planetIDs {
		if pid.IsCapitol {
			capitolID = pid.PlanetID
			break
		}
	}

	planet, err := s.getPlanetByID(ctx, capitolID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("getPlanetByID(): %w", err)
	}

	return planet, nil
}
