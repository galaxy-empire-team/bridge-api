package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetPlanet(ctx context.Context, planetID uuid.UUID) (models.Planet, error) {
	planet, err := s.getPlanetByID(ctx, planetID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("getPlanetByID(): %w", err)
	}

	return planet, nil
}
