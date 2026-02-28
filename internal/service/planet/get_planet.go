package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetPlanet(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.Planet, error) {
	isUserPlanet, err := s.planetStorage.CheckPlanetBelongsToUser(ctx, userID, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetStorage.CheckPlanetBelongsToUser(): %w", err)
	}
	if !isUserPlanet {
		return models.Planet{}, models.ErrPlanetDoesNotBelongToUser
	}

	planet, err := s.getPlanetByID(ctx, userID, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("getPlanetByID(): %w", err)
	}

	return planet, nil
}
