package rating

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetRating(ctx context.Context, userID uuid.UUID) (ratings models.Ratings, err error) {
	errG, gCtx := errgroup.WithContext(ctx)
	errG.Go(func() error {
		userRating, err := s.getUsersRating(gCtx, userID)
		if err != nil {
			return fmt.Errorf("getUsersRating(): %w", err)
		}

		ratings.User = userRating

		return nil
	})

	errG.Go(func() error {
		fleetRating, err := s.getFleetRating(gCtx, userID)
		if err != nil {
			return fmt.Errorf("getFleetRating(): %w", err)
		}

		ratings.Fleet = fleetRating

		return nil
	})

	errG.Go(func() error {
		playersCount, err := s.userStorage.GetUserCount(gCtx)
		if err != nil {
			return fmt.Errorf("userStorage.GetUserCount(): %w", err)
		}

		ratings.PlayersCount = playersCount

		return nil
	})

	if err := errG.Wait(); err != nil {
		return models.Ratings{}, fmt.Errorf("errG.Wait(): %w", err)
	}

	return ratings, nil
}
