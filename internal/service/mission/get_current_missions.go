package mission

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetCurrentMissions(ctx context.Context, userID uuid.UUID) (missions models.UserMissions, err error) {
	errG, ctxG := errgroup.WithContext(ctx)

	errG.Go(func() error {
		missions.Missions, err = s.missionStorage.GetCurrentUserMissions(ctxG, userID)
		if err != nil {
			return fmt.Errorf("missionStorage.GetCurrentUserMissions(): %w", err)
		}

		return nil
	})

	errG.Go(func() error {
		missions.MistLaunches, err = s.planetStorage.GetMistLaunches(ctxG, userID)
		if err != nil {
			return fmt.Errorf("planetStorage.GetMistLaunches(): %w", err)
		}

		return nil
	})

	return missions, errG.Wait()
}
