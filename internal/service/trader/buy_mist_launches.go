package trader

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func (s *Service) buyMistLaunches(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, itemStats registry.TraderItemStats) error {
	mistLaunchItem, ok := itemStats.Item.(registry.TraderItemAutoMistLaunch)
	if !ok {
		return fmt.Errorf("typecast to registry.TraderItemAutoMistLaunch")
	}

	return s.txManager.ExecTraderTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.removeResources(ctx, userID, planetID, itemStats.Cost, storages)
		if err != nil {
			return fmt.Errorf("removeResources(): %w", err)
		}

		err = storages.AddMistLaunches(ctx, userID, mistLaunchItem.Count)
		if err != nil {
			return fmt.Errorf("AddMistLaunches(): %w", err)
		}

		return nil
	})
}
