package research

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *ResearchStorage) GetUserResearchesProgress(ctx context.Context, userID uuid.UUID) ([]models.ResearchProgressInfo, error) {
	const getAllUserResearchesQuery = `
		SELECT 
			user_id,
			research_id,
			started_at,
			finished_at
		FROM session_beta.event_researches
		WHERE user_id = $1;
	`

	rows, err := r.DB.Query(ctx, getAllUserResearchesQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	var researches []models.ResearchProgressInfo
	for rows.Next() {
		var research models.ResearchProgressInfo
		var startedAt, finishedAt time.Time
		err = rows.Scan(
			&research.UserID,
			&research.ResearchID,
			&startedAt,
			&finishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		research.StartedAt = startedAt.UTC()
		research.FinishedAt = finishedAt.UTC()

		researches = append(researches, research)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return researches, nil
}
