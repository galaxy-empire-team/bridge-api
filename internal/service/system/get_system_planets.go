package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetSystemPlanets(ctx context.Context, userID uuid.UUID, system models.System) (models.SystemPlanets, error) {
	var (
		errGroup errgroup.Group
		planets  models.SystemPlanets
		err      error
	)

	errGroup.Go(func() error {
		planets.Planets, err = s.systemStorage.GetSystemPlanets(ctx, system)
		if err != nil {
			return fmt.Errorf("systemStorage.GetSystemPlanets(): %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		planets.NPC, err = s.planetStorage.GetUserNPCAttacks(ctx, userID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetUserNPCAttacks(): %w", err)
		}

		return nil
	})

	if err = errGroup.Wait(); err != nil {
		return models.SystemPlanets{}, fmt.Errorf("errGroup.Wait(): %w", err)
	}

	planets.System = system

	return planets, nil
}
