package rating

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetRating(ctx context.Context, userID uuid.UUID) (models.Ratings, error) {
	var (
		errG    errgroup.Group
		ratings models.Ratings
	)

	errG.Go(func() error {
		userRating, err := s.getUsersRating(ctx, userID)
		if err != nil {
			return fmt.Errorf("getUsersRating(): %w", err)
		}

		ratings.User = userRating

		return nil
	})

	errG.Go(func() error {
		fleetRating, err := s.getFleetRating(ctx, userID)
		if err != nil {
			return fmt.Errorf("getFleetRating(): %w", err)
		}

		ratings.Fleet = fleetRating

		return nil
	})

	if err := errG.Wait(); err != nil {
		return models.Ratings{}, fmt.Errorf("errG.Wait(): %w", err)
	}

	return ratings, nil
}
