package event

import (
	"context"
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (s *EventStorage) SetResearchFinishTime(ctx context.Context, researchEvent models.EventFinishTime) error {
	const setResearchFinishTimeQuery = `
		UPDATE session_beta.event_researches
		SET finished_at = $1
		WHERE id = $2;
	`

	_, err := s.DB.Exec(ctx, setResearchFinishTimeQuery, researchEvent.FinishedAt, researchEvent.EventID)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
