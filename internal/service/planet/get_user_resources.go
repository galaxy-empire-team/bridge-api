package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetUserResources(ctx context.Context, userID uuid.UUID) (res GetUserResourcesResponse, err error) {
	errG, gCtx := errgroup.WithContext(ctx)
	errG.Go(func() error {
		res.UserResources, err = s.planetStorage.GetUserResources(gCtx, userID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetUserResources(): %w", err)
		}

		return nil
	})

	errG.Go(func() error {
		res.Boosts, err = s.planetStorage.GetUserBoosts(gCtx, userID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetUserBoosts(): %w", err)
		}

		return nil
	})

	if err := errG.Wait(); err != nil {
		return GetUserResourcesResponse{}, fmt.Errorf("errG.Wait(): %w", err)
	}

	return res, nil
}
