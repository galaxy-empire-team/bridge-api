package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

func (s *Service) CreateCapitol(ctx context.Context, userID uuid.UUID) error {
	planetID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("uuid.NewV7(): %w", err)
	}

	generatedLocation := models.Location{
		X: uint8(s.randomGenerator.Uint32() % galaxyCount),
		Y: uint8(s.randomGenerator.Uint32() % systemInGalaxyCount),
		Z: uint8(s.randomGenerator.Uint32() % planetsInSystemCount),
	}
	planetToColonize := models.Planet{
		ID:        planetID,
		Location:  generatedLocation,
		IsCapitol: true,
		HasMoon:   false,
		Buildings: make(map[models.BuildingType]models.BuildingInfo),
	}

	for _, buildingType := range models.GetAllBuildings() {
		planetToColonize.Buildings[buildingType] = models.BuildingInfo{
			Type:  buildingType,
			Level: defaultLvl,
		}
	}
	err = s.planetStorage.CreatePlanet(ctx, userID, planetToColonize)
	if err != nil {
		return fmt.Errorf("planetRepo.CreatePlanet(): %w", err)
	}

	return nil
}
