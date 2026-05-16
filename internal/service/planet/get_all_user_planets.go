package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetAllUserPlanets(ctx context.Context, userID uuid.UUID) ([]models.Planet, error) {
	planetIDs, err := s.planetStorage.GetUserPlanetIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("planetRepo.GetUserPlanetIDs(): %w", err)
	}

	if len(planetIDs) == 0 {
		return nil, models.ErrNoPlanetsFound
	}

	userPlanets, err := s.planetStorage.GetAllUserPlanets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("planetStorage.GetAllUserPlanets(): %w", err)
	}

	return userPlanets, nil
}
