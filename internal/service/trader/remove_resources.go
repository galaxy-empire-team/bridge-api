package trader

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func (s *Service) removeResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, cost registry.Resources, storages TxStorages) error {
	err := s.removePlanetResources(ctx, planetID, models.Resources{
		Metal:   cost.Metal,
		Crystal: cost.Crystal,
		Gas:     cost.Gas,
	}, storages)
	if err != nil {
		return fmt.Errorf("removePlanetResources(): %w", err)
	}

	err = s.removeUserResources(ctx, userID, models.UserResources{
		Matter: cost.Matter,
		Doreye: cost.Doreye,
	}, storages)
	if err != nil {
		return fmt.Errorf("removeUserResources(): %w", err)
	}

	return nil
}

func (s *Service) removePlanetResources(ctx context.Context, planetID uuid.UUID, cost models.Resources, storages TxStorages) error {
	if cost.IsEmpty() {
		return nil
	}

	planetResources, err := storages.GetResourcesForUpdate(ctx, planetID)
	if err != nil {
		return fmt.Errorf("GetResourcesForUpdate(): %w", err)
	}

	if planetResources.Metal < cost.Metal || planetResources.Crystal < cost.Crystal || planetResources.Gas < cost.Gas {
		return models.ErrNotEnoughResources
	}

	updatedResources := models.Resources{
		Metal:     planetResources.Metal - cost.Metal,
		Crystal:   planetResources.Crystal - cost.Crystal,
		Gas:       planetResources.Gas - cost.Gas,
		UpdatedAt: planetResources.UpdatedAt,
	}

	err = storages.SetResources(ctx, planetID, updatedResources)
	if err != nil {
		return fmt.Errorf("SetResources(): %w", err)
	}

	return nil
}

func (s *Service) removeUserResources(ctx context.Context, userID uuid.UUID, cost models.UserResources, storages TxStorages) error {
	if cost.IsEmpty() {
		return nil
	}

	userResources, err := storages.GetUserResourcesForUpdate(ctx, userID)
	if err != nil {
		return fmt.Errorf("storages.GetUserResourcesForUpdate(): %w", err)
	}

	if userResources.Matter < cost.Matter {
		return models.ErrNotEnoughMatter
	}

	if userResources.Doreye < cost.Doreye {
		return models.ErrNotEnoughDoreye
	}

	updatedUserResources := models.UserResources{
		Matter: userResources.Matter - cost.Matter,
		Doreye: userResources.Doreye - cost.Doreye,
	}

	err = storages.SetUserResources(ctx, userID, updatedUserResources)
	if err != nil {
		return fmt.Errorf("SetUserResources(): %w", err)
	}

	return nil
}
