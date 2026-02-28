package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) updateResources(ctx context.Context, planetID uuid.UUID, cargo models.Resources, storage TxStorages) error {
	resources, err := storage.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return fmt.Errorf("planetStorage.GetResourcesForUpdate(): %w", err)
	}

	if resources.Metal < cargo.Metal || resources.Crystal < cargo.Crystal || resources.Gas < cargo.Gas {
		return models.ErrNotEnoughResources
	}

	updatedResources := models.Resources{
		Metal:     resources.Metal - cargo.Metal,
		Crystal:   resources.Crystal - cargo.Crystal,
		Gas:       resources.Gas - cargo.Gas,
		UpdatedAt: resources.UpdatedAt,
	}

	err = storage.SetResources(ctx, planetID, updatedResources)
	if err != nil {
		return fmt.Errorf("planetStorage.SetResources(): %w", err)
	}

	return nil
}
