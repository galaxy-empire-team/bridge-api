package event

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) useBoost(ctx context.Context, userID uuid.UUID, boost models.UserBoost, txStorages TxStorages) (time.Duration, error) {
	userBoostCount, err := txStorages.GetBoostByIDForUpdate(ctx, userID, boost.ID)
	if err != nil {
		return 0, fmt.Errorf("txStorages.GetBoostByIDForUpdate(): %w", err)
	}

	boostStats, err := s.registry.GetBoostStatsByID(boost.ID)
	if err != nil {
		return 0, fmt.Errorf("registry.GetBoostStatsByID(): %w", err)
	}

	if userBoostCount.Count < boost.Count {
		return 0, models.ErrNotEnoughBoosts
	}

	userBoostCount.Count -= boost.Count
	err = txStorages.SetBoost(ctx, userID, userBoostCount)
	if err != nil {
		return 0, fmt.Errorf("txStorages.SetBoost(): %w", err)
	}

	duration := time.Duration(boostStats.DurationS*boost.Count) * time.Second

	return duration, nil
}
