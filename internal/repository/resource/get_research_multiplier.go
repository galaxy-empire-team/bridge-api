package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Repository) getResearchResourceMultiplier(ctx context.Context, userID uuid.UUID) (float32, error) {
	userResearches, err := s.researchStorage.GetUserResearchesByTypes(ctx, userID, []consts.ResearchType{consts.ResearchTypeIndustrialTechnology})
	if err != nil {
		return 0, fmt.Errorf("researchStorage.GetUserResearchesByTypes(): %w", err)
	}

	researchID, ok := userResearches[consts.ResearchTypeIndustrialTechnology]
	if !ok {
		return consts.BaseResourceGenerationSpeed, nil
	}

	researchStats, err := s.registry.GetResearchStatsByID(researchID)
	if err != nil {
		return 0, fmt.Errorf("registry.GetResearchStatsByID(): %w", err)
	}

	return researchStats.Bonuses.ProductionSpeedImprove, nil
}
