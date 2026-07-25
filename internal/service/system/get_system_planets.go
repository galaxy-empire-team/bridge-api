package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetSystemPlanets(ctx context.Context, userID uuid.UUID, system models.System) (planets models.SystemPlanets, err error) {
	errG, gCtx := errgroup.WithContext(ctx)
	errG.Go(func() error {
		planets.Planets, err = s.systemStorage.GetSystemPlanets(gCtx, system)
		if err != nil {
			return fmt.Errorf("systemStorage.GetSystemPlanets(): %w", err)
		}

		return nil
	})

	errG.Go(func() error {
		planets.NPC, err = s.planetStorage.GetUserNPCAttacks(gCtx, userID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetUserNPCAttacks(): %w", err)
		}

		return nil
	})

	if err = errG.Wait(); err != nil {
		return models.SystemPlanets{}, fmt.Errorf("errG.Wait(): %w", err)
	}

	planets.System = system

	return planets, nil
}
