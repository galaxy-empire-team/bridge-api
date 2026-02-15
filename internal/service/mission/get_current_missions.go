package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetCurrentMissions(ctx context.Context, userID uuid.UUID) ([]models.UserMission, error) {
	missions, err := s.missionStorage.GetCurrentUserMissions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("missionStorage.GetCurrentUserMissions(): %w", err)
	}

	return missions, nil
}
