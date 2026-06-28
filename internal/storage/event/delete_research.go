package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *EventStorage) DeleteResearchEvent(ctx context.Context, userID uuid.UUID, researchID consts.ResearchID) error {
	const deleteResearchEventQuery = `
			DELETE FROM session_beta.event_researches
			WHERE user_id = $1 AND research_id = $2;
		`

	cmd, err := s.DB.Exec(ctx, deleteResearchEventQuery, userID, researchID)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrEventIsNotScheduled
	}

	return nil
}
