package research

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *ResearchStorage) CreateResearchEvent(ctx context.Context, researchEvent models.ResearchEvent) error {
	const createResearchEventQuery = `
		INSERT INTO session_beta.event_researches (
			user_id,
			research_id, 
			started_at,
			finished_at
		) VALUES (
			$1,    -- user_id
			$2,    -- research_id
			$3,    -- started_at
			$4	   -- finished_at
		) ON CONFLICT (user_id) DO NOTHING;
		`

	cmd, err := s.DB.Exec(ctx, createResearchEventQuery,
		researchEvent.UserID,
		researchEvent.ResearchID,
		researchEvent.StartedAt,
		researchEvent.FinishedAt,
	)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrEventIsAlreadyScheduled
	}

	return nil
}
