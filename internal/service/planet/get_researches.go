package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *Service) GetResearches(ctx context.Context, userID uuid.UUID) (models.UserResearches, error) {
	researches, err := s.researchStorage.GetUserResearches(ctx, userID)
	if err != nil {
		return models.UserResearches{}, fmt.Errorf("researchStorage.GetUserResearches(): %w", err)
	}

	researchProgress, err := s.researchStorage.GetUserResearchesProgress(ctx, userID)
	if err != nil {
		return models.UserResearches{}, fmt.Errorf("researchStorage.GetUserResearchesProgress(): %w", err)
	}

	return models.UserResearches{
		UserID:           userID,
		Research:         researches,
		ResearchProgress: researchProgress,
	}, nil
}
