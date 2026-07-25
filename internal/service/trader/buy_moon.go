package trader

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func (s *Service) buyMoon(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, itemStats registry.TraderItemStats) error {
	_, ok := itemStats.Item.(registry.TraderItemMoon)
	if !ok {
		return fmt.Errorf("typecast to registry.TraderItemMoon")
	}

	return s.txManager.ExecTraderTx(ctx, func(ctx context.Context, storages TxStorages) error {
		hasMoon, err := storages.CheckPlanetHasMoon(ctx, planetID)
		if err != nil {
			return fmt.Errorf("CheckPlanetHasMoon(): %w", err)
		}
		if hasMoon {
			return models.ErrMoonAlreadyExists
		}

		err = s.removeResources(ctx, userID, planetID, itemStats.Cost, storages)
		if err != nil {
			return fmt.Errorf("removeResources(): %w", err)
		}

		err = storages.SetHasMoon(ctx, planetID, true)
		if err != nil {
			return fmt.Errorf("SetHasMoon(): %w", err)
		}

		return nil
	})
}
