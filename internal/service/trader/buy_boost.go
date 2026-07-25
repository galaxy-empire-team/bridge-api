package trader

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func (s *Service) buyBoost(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, itemStats registry.TraderItemStats) error {
	boostItem, ok := itemStats.Item.(registry.TraderItemBoost)
	if !ok {
		return fmt.Errorf("typecast to registry.TraderItemBoost")
	}

	return s.txManager.ExecTraderTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.removeResources(ctx, userID, planetID, itemStats.Cost, storages)
		if err != nil {
			return fmt.Errorf("removeResources(): %w", err)
		}

		err = storages.AddBoost(ctx, userID, models.UserBoost{
			ID:    boostItem.ID,
			Count: boostItem.Count,
		})
		if err != nil {
			return fmt.Errorf("AddBoost(): %w", err)
		}

		return nil
	})
}
