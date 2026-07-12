package planet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) ActivateMoon(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, boostID consts.MoonBoostID, count uint64) (models.MoonInfo, error) {
	if count == 0 {
		return models.MoonInfo{}, models.ErrInvalidInput
	}

	if err := s.repository.CheckPlanetOwner(ctx, userID, planetID); err != nil {
		return models.MoonInfo{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	hasMoon, err := s.planetStorage.CheckPlanetHasMoon(ctx, planetID)
	if err != nil {
		return models.MoonInfo{}, fmt.Errorf("CheckPlanetHasMoon(): %w", err)
	}
	if !hasMoon {
		return models.MoonInfo{}, models.ErrMoonNotFound
	}

	boostStats, err := s.registry.GetMoonBoostStatsByID(boostID)
	if err != nil {
		return models.MoonInfo{}, fmt.Errorf("registry.GetMoonBoostStatsByID(): %w", err)
	}

	matterCost := boostStats.MatterCost * count
	duration := time.Duration(boostStats.DurationS*count) * time.Second

	moonInfo := models.MoonInfo{
		PlanetID: planetID,
		HasMoon:  true,
	}

	return moonInfo, s.txManager.ExecPlanetTx(ctx, func(ctx context.Context, storageTx TxStorages) error {
		userResources, err := storageTx.GetUserResourcesForUpdate(ctx, userID)
		if err != nil {
			return fmt.Errorf("GetUserResourcesForUpdate(): %w", err)
		}

		if userResources.Matter < matterCost {
			return models.ErrNotEnoughMatter
		}

		err = storageTx.SetUserMatter(ctx, userID, userResources.Matter-matterCost)
		if err != nil {
			return fmt.Errorf("SetUserMatter(): %w", err)
		}

		currentMoonInfo, err := storageTx.GetMoonActivationForUpdate(ctx, planetID)
		if err != nil {
			return fmt.Errorf("GetMoonActivationForUpdate(): %w", err)
		}

		now := time.Now().UTC()
		baseTime := now
		if currentMoonInfo.ActivateUntill.After(now) {
			baseTime = currentMoonInfo.ActivateUntill
		}

		activateUntill := baseTime.Add(duration)
		err = storageTx.SetMoonActivation(ctx, planetID, activateUntill)
		if err != nil {
			return fmt.Errorf("SetMoonActivation(): %w", err)
		}

		moonInfo.ActivateUntill = activateUntill

		return nil
	})
}
