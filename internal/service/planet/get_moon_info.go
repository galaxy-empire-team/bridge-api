package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetMoonInfo(ctx context.Context, userID uuid.UUID, planetID uuid.UUID) (models.MoonInfo, error) {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return models.MoonInfo{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	hasMoon, err := s.planetStorage.CheckPlanetHasMoon(ctx, planetID)
	if err != nil {
		return models.MoonInfo{}, fmt.Errorf("CheckPlanetHasMoon(): %w", err)
	}
	if !hasMoon {
		return models.MoonInfo{
			PlanetID: planetID,
			HasMoon:  false,
		}, nil
	}

	moonInfo, err := s.planetStorage.GetMoonActivationForUpdate(ctx, planetID)
	if err != nil {
		return models.MoonInfo{}, fmt.Errorf("GetMoonActivationForUpdate(): %w", err)
	}

	return models.MoonInfo{
		PlanetID:       planetID,
		HasMoon:        true,
		ActivateUntill: moonInfo.ActivateUntill,
	}, nil
}
