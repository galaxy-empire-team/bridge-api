package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *EventStorage) GetResearchEventForUpdate(ctx context.Context, userID uuid.UUID, researchID consts.ResearchID) (models.EventFinishTime, error) {
	const getResearchEventForUpdateQuery = `
			SELECT 
				id,
				started_at,
				finished_at
			FROM session_beta.event_researches
			WHERE user_id = $1 AND research_id = $2
			FOR UPDATE;
		`

	var researchEvent models.EventFinishTime
	err := s.DB.QueryRow(ctx, getResearchEventForUpdateQuery, userID, researchID).Scan(&researchEvent.EventID, &researchEvent.StartedAt, &researchEvent.FinishedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.EventFinishTime{}, models.ErrEventIsNotScheduled
		}
		return models.EventFinishTime{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return researchEvent, nil
}
