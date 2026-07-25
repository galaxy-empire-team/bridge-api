package trader

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Buy(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, itemID consts.TraderItemID) error {
	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	itemStats, err := s.registry.GetTraderItemStatsByID(itemID)
	if err != nil {
		return fmt.Errorf("registry.GetTraderItemStatsByID(): %w", err)
	}

	if itemStats.Cost.Metal > 0 || itemStats.Cost.Crystal > 0 || itemStats.Cost.Gas > 0 {
		err := s.repository.RecalcResources(ctx, userID, planetID)
		if err != nil {
			return fmt.Errorf("RecalcResources(): %w", err)
		}
	}

	switch itemStats.ItemType {
	case consts.TraderItemTypeBoost:
		err = s.buyBoost(ctx, userID, planetID, itemStats)
		if err != nil {
			return fmt.Errorf("buyBoost(): %w", err)
		}
	case consts.TraderItemTypeMoon:
		err = s.buyMoon(ctx, userID, planetID, itemStats)
		if err != nil {
			return fmt.Errorf("buyMoon(): %w", err)
		}
	case consts.TraderItemTypeSpaceship:
		err = s.buySpaceship(ctx, userID, planetID, itemStats)
		if err != nil {
			return fmt.Errorf("buySpaceship(): %w", err)
		}
	case consts.TraderItemTypeAutoMistLaunch:
		err = s.buyMistLaunches(ctx, userID, planetID, itemStats)
		if err != nil {
			return fmt.Errorf("buyMistLaunches(): %w", err)
		}
	default:
		return models.ErrUnknownTraderItemType
	}

	return nil
}
