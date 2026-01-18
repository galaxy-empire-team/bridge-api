package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (s *Service) GetCapitolPlanet(ctx context.Context, userID uuid.UUID) (models.Planet, error) {
	capitolPlanet, err := s.planetRepo.GetCapitol(ctx, userID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetCapitol(): %w", err)
	}

	return capitolPlanet, nil
}
