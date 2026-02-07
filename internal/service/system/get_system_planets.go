package planet

import (
	"context"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetSystemPlanets(ctx context.Context, system models.System) (models.SystemPlanets, error) {
	planets, err := s.systemStorage.GetSystemPlanets(ctx, system)
	if err != nil {
		return models.SystemPlanets{}, err
	}

	return planets, nil
}
