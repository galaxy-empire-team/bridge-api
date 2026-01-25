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

	if capitolID == uuid.Nil {
		return models.Planet{}, models.ErrCapitolNotFound
	}

	err = s.recalcResources(ctx, capitolID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("recalcResources(): %w", err)
	}

	capitolLocation, err := s.planetStorage.GetLocation(ctx, capitolID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetCapitol(): %w", err)
	}

	resources, err := s.planetStorage.GetResources(ctx, capitolID)
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetResources(): %w", err)
	}

	buildings, err := s.planetStorage.GetBuildingsInfo(ctx, capitolID, models.GetAllBuildings())
	if err != nil {
		return models.Planet{}, fmt.Errorf("planetRepo.GetBuildingsInfo(): %w", err)
	}

	return models.Planet{
		ID:        capitolID,
		Location:  capitolLocation,
		Resources: resources,
		Buildings: buildings,
		IsCapitol: true,
		HasMoon:   false,
	}, nil
}
