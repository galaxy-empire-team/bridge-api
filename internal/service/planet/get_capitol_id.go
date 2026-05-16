package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) GetCapitolID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	planet, err := s.planetStorage.GetCapitolID(ctx, userID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("getPlanetByID(): %w", err)
	}

	return planet, nil
}
