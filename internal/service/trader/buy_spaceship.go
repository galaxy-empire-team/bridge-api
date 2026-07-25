package trader

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/registry"
)

func (s *Service) buySpaceship(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, itemStats registry.TraderItemStats) error {
	spaceshipItem, ok := itemStats.Item.(registry.TraderItemSpaceship)
	if !ok {
		return fmt.Errorf("invalid trader item cast to TraderItemSpaceship")
	}

	return s.txManager.ExecTraderTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.removeResources(ctx, userID, planetID, itemStats.Cost, storages)
		if err != nil {
			return fmt.Errorf("removeResources(): %w", err)
		}

		err = storages.AddFleet(ctx, planetID, models.FleetUnitCount{
			ID:    spaceshipItem.ID,
			Count: spaceshipItem.Count,
		})
		if err != nil {
			return fmt.Errorf("AddFleet(): %w", err)
		}

		return nil
	})
}
