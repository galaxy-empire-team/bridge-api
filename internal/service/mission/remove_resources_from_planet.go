package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) removeResourcesFromPlanet(ctx context.Context, planetID uuid.UUID, amount models.Resources, storage TxStorages) error {
	resources, err := storage.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetResourcesForUpdate(): %w", err)
	}

	if resources.Metal < amount.Metal || resources.Crystal < amount.Crystal || resources.Gas < amount.Gas {
		return models.ErrNotEnoughResources
	}

	updatedResources := models.Resources{
		Metal:     resources.Metal - amount.Metal,
		Crystal:   resources.Crystal - amount.Crystal,
		Gas:       resources.Gas - amount.Gas,
		UpdatedAt: resources.UpdatedAt,
	}

	err = storage.SetResources(ctx, planetID, updatedResources)
	if err != nil {
		return fmt.Errorf("planetStorage.SetResources(): %w", err)
	}

	return nil
}
